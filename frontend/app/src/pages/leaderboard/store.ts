/* eslint-disable prefer-destructuring */
import { orderBy } from 'lodash';
import memo from 'memo-decorator';
import { makeAutoObservable } from 'mobx';
import api from '../../api';

export type LeaderItem = {
  owner_pubkey: string;
  total_bounties_completed: number;
  total_sats_earned: number;
};

export class LeaderboardStore {
  private leaders: LeaderItem[] = [];

  public isLoading = false;
  public error: any;
  public total: LeaderItem | null = null;
  constructor() {
    makeAutoObservable(this);
  }

  @memo()
  async fetchLeaders() {
    this.isLoading = true;
    try {
      const resp = (await api.get('people/bounty/leaderboard')) as LeaderItem[];
      this.total = resp.reduce(
        (partialSum: LeaderItem, assigneeStats: LeaderItem) => {
          partialSum.total_bounties_completed += assigneeStats.total_bounties_completed;
          partialSum.total_sats_earned += assigneeStats.total_sats_earned;

          return partialSum;
        },
        { owner_pubkey: '', total_bounties_completed: 0, total_sats_earned: 0 }
      );
      this.leaders = resp;
    } catch (e) {
      this.error = e;
    } finally {
      this.isLoading = false;
    }
  }

  get sortedBySats() {
    return orderBy(this.leaders, 'total_sats_earned', 'desc');
  }

  get top3() {
    return this.sortedBySats.slice(0, 3);
  }
  get others() {
    return this.sortedBySats.slice(3);
  }
}

export const leaderboardStore = new LeaderboardStore();
