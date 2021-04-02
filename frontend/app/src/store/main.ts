import { observable, action } from 'mobx'
import api from '../api'

export class MainStore {
  @observable
  tribes: Tribe[] = []

  @action async getTribes(uniqueName?:string): Promise<Tribe[]> {
    const ts = await api.get('tribes')
    ts.sort((a:Tribe,b:Tribe)=>{
      if (b.last_active===a.last_active) {
        return b.member_count-a.member_count
      }
      return b.last_active-a.last_active
    })
    if(uniqueName) {
      ts.forEach(function(t:Tribe,i:number){
        if(t.unique_name === uniqueName){
          ts.splice(i, 1);
          ts.unshift(t);
        }
      })
    }
    this.tribes = ts
    return ts
  }
}

export const mainStore = new MainStore()

export interface Tribe {
  uuid: string;
  name: string;
  unique_name: string;
  owner: string;
  pubkey: string; // group encryption key
  price: number;
  img: string;
  tags: string[];
  description: string;
  member_count: number;
  last_active: number;
  matchCount?: number; // for tag search
}

