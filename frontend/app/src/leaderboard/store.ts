/* eslint-disable prefer-destructuring */
import api from 'api';
import { orderBy } from 'lodash';
import memo from 'memo-decorator';
import { makeAutoObservable } from 'mobx';

export type LeaderItem = {
  owner_pubkey: string;
  total_bounties_completed: number;
  total_sats_earned: number;
};
export class LeaderboardStore {
  private leaders: LeaderItem[] = [];

  public isLoading = false;
  public error: any;
  public total: Omit<LeaderItem, 'owner_pubkey'> | null = null;
  constructor() {
    makeAutoObservable(this);
  }

  @memo()
  async fetchLeaders() {
    this.isLoading = true;
    try {
      const resp = (await api.get('people/bounty/leaderboard')) as LeaderItem[];
      this.total = resp[0];
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
