import { observable, action } from 'mobx'
import api from '../api'

export class MainStore {
  @observable
  tribes: Tribe[] = []

  @action async getTribes(){
    const ts = await api.get('tribes')
    this.tribes = ts
  }
}

export const mainStore = new MainStore()

export interface Tribe {
  uuid: string;
  name: string;
  owner: string;
  pubkey: string; // group encryption key
  price: number;
  img: string;
  tags: string[];
  description: string;
  member_count: number;
  matchCount?: number; // for tag search
}

