import { observable, action } from 'mobx'
import { persist } from 'mobx-persist'
import api from '../api'
import { Extras } from '../form/inputs/widgets/interfaces'

export class MainStore {
  @persist('list') @observable
  tribes: Tribe[] = []

  @action async getTribes(uniqueName?: string): Promise<Tribe[]> {
    const ts = await api.get('tribes')
    ts.sort((a: Tribe, b: Tribe) => {
      if (b.last_active === a.last_active) {
        return b.member_count - a.member_count
      }
      return b.last_active - a.last_active
    })
    if (uniqueName) {
      ts.forEach(function (t: Tribe, i: number) {
        if (t.unique_name === uniqueName) {
          ts.splice(i, 1);
          ts.unshift(t);
        }
      })
    }
    this.tribes = ts
    return ts
  }

  @persist('list') @observable
  people: Person[] = []

  @action async getPeople(uniqueName?: string): Promise<Person[]> {
    const ps = await api.get('people')
    if (uniqueName) {
      ps.forEach(function (t: Tribe, i: number) {
        if (t.unique_name === uniqueName) {
          ps.splice(i, 1);
          ps.unshift(t);
        }
      })
    }
    this.people = ps
    return ps
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

export interface Person {
  id: number;
  unique_name: string;
  owner_pubkey: string;
  owner_alias: string;
  description: string;
  img: string;
  tags: string[];
  pubkey: string
  photo_url: string
  alias: string
  route_hint: string
  contact_key: string
  price_to_meet: number
  url: string
  verification_signature: string
  extras: Extras
}
