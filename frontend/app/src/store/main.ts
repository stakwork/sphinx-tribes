import { observable, action } from 'mobx'
import { persist } from 'mobx-persist'
import api from '../api'
import { Extras } from '../form/inputs/widgets/interfaces'
import { getHostIncludingDockerHosts } from '../host'
import { uiStore } from './ui'

export class MainStore {
  @persist('list') @observable
  tribes: Tribe[] = []
  ownerTribes: Tribe[] = []

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

  bots: Bot[] = []

  @action async getBots(uniqueName?: string): Promise<Bot[]> {
    let b = await api.get('bots')

    if (uniqueName) {
      b.forEach(function (t: Bot, i: number) {
        if (t.unique_name === uniqueName) {
          b.splice(i, 1);
          b.unshift(t);
        }
      })
    }

    // b = [{
    //   name: 'welcome',
    //   unique_name: 'welcome',
    //   label: 'Welcome',
    //   description: 'my first bot bot'
    // }, {
    //   name: 'btc',
    //   unique_name: 'btc',
    //   label: 'BTC',
    //   description: 'my first bot bot'
    // }, {
    //   name: 'bet',
    //   unique_name: 'bet',
    //   label: 'Bet',
    //   description: 'my first bot botmy first bot botmy first bot botmy first bot bot'
    // },]

    // hide test bots and set images
    b && b.forEach((bb, i) => {
      if (bb.unique_name === 'btc') {
        bb.img = '/static/bots_bitcoin.png'
        b.splice(i, 1);
        b.unshift(bb);
      }
      if (bb.unique_name === 'bet') {
        bb.img = '/static/bots_betting.png'
        b.splice(i, 1);
        b.unshift(bb);
      }
      if (bb.unique_name === 'hello' || bb.unique_name === 'welcome') {
        bb.img = '/static/bots_welcome.png'
        b.splice(i, 1);
        b.unshift(bb);
      }
      if (bb.unique_name && bb.unique_name.includes('test')) {
        // hide all test bots
        bb.hide = true
      }
    })


    this.bots = b
    return b
  }

  @action async getTribesByOwner(pubkey: string): Promise<Tribe[]> {
    const ts = await api.get(`tribes_by_owner/${pubkey}`)
    this.ownerTribes = ts
    return ts
  }

  @action async makeBot(payload: any): Promise<Bot> {
    const b = await api.post('bots', payload)
    console.log('made bot', b)
    return b
  }

  @persist('list') @observable
  people: Person[] = []

  @action async getPeople(uniqueName?: string): Promise<Person[]> {
    const ps = await api.get('people')

    if (uiStore.meInfo) {
      const index = ps.findIndex(f => f.id == uiStore.meInfo?.id)

      if (index > -1) {
        // add 'hide' property to me in people list
        ps[index].hide = true

        if (!uiStore.meInfo.img && ps[index].img) {
          // if meInfo has no img but people list does, copy that image to meInfo
          uiStore.setMeInfo({ ...uiStore.meInfo, img: ps[index].img })
        }
      }
    }
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

  @action async refreshJwt() {
    try {
      if (!uiStore.meInfo) return null
      const info = uiStore.meInfo
      const URL = info.url.startsWith("http") ? info.url : `https://${info.url}`;
      const res: any = await fetch(URL + "/refresh_jwt", {
        method: "GET",
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });
      const j = await res.json()

      return j.response
    } catch (e) {
      console.log('e', e)
      // could not refresh jwt, logout!
      return null
    }
  }

  @action async deleteProfile() {
    try {
      if (!uiStore.meInfo) return null
      const info = uiStore.meInfo
      const URL = info.url.startsWith("http") ? info.url : `https://${info.url}`;
      const res: any = await fetch(URL + "/profile", {
        method: "DELETE",
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...info
        }),
      });

      uiStore.setMeInfo(null)
      uiStore.setSelectingPerson(0)
      uiStore.setSelectedPerson(0)

      const j = await res.json()
      return j
    } catch (e) {
      console.log('e', e)
      // could not delete profile!
      return null
    }
  }

  @action async addFavorite() {
    let body: any = {}
    console.log('SUBMIT FORM', body);

    // console.log('mergeFormWithMeData', body);
    if (!body) return // avoid saving bad state

    const info = uiStore.meInfo as any;
    if (!info) return console.log("no meInfo");
    try {
      const URL = info.url.startsWith("http") ? info.url : `https://${info.url}`;
      const r = await fetch(URL + "/profile", {
        method: "POST",
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...body,
          price_to_meet: parseInt(body.price_to_meet),
        }),
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      })


      if (!r.ok) {
        return alert("Failed to save data");
      }

      uiStore.setToasts([{
        id: '1',
        title: 'Added to favorites.'
      }]);

    } catch (e) {
      console.log('e', e)
    }
  }

  @action async deleteFavorite() {
    let body: any = {}
    console.log('SUBMIT FORM', body);

    // console.log('mergeFormWithMeData', body);
    if (!body) return // avoid saving bad state

    const info = uiStore.meInfo as any;
    if (!info) return console.log("no meInfo");
    try {
      const URL = info.url.startsWith("http") ? info.url : `https://${info.url}`;
      const r = await fetch(URL + "/profile", {
        method: "POST",
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...body,
          price_to_meet: parseInt(body.price_to_meet),
        }),
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      })


      if (!r.ok) {
        return alert("Failed to save data");
      }

      uiStore.setToasts([{
        id: '1',
        title: 'Added to favorites.'
      }]);

    } catch (e) {
      console.log('e', e)
    }
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

export interface Bot {
  id?: number;
  uuid: string;
  name: string;
  owner_pubkey: string;
  unique_name: string;
  price_per_use: number;
  created: string;
  updated: string;
  unlisted: boolean;
  deleted: boolean;
  owner_route_hint: string;
  owner: string;
  pubkey: string; // group encryption key
  price: number;
  img: string;
  tags: string[];
  description: string;
  member_count: number;
  hide?: boolean;
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
  hide?: boolean
}

export interface Jwt {
  jwt: string;
}