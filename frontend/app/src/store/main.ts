import { observable, action } from "mobx";
import { persist } from "mobx-persist";
import api from "../api";
import { Extras } from "../form/inputs/widgets/interfaces";
import { getHostIncludingDockerHosts } from "../host";
import { MeInfo, uiStore } from "./ui";

export const queryLimit = 100
export const smallQueryLimit = 100

export class MainStore {
  @persist("list")
  @observable
  tribes: Tribe[] = [];
  ownerTribes: Tribe[] = [];

  @action async getTribes(uniqueName?: string): Promise<Tribe[]> {
    const ts = await api.get("tribes");
    ts.sort((a: Tribe, b: Tribe) => {
      if (b.last_active === a.last_active) {
        return b.member_count - a.member_count;
      }
      return b.last_active - a.last_active;
    });
    if (uniqueName) {
      ts.forEach(function (t: Tribe, i: number) {
        if (t.unique_name === uniqueName) {
          ts.splice(i, 1);
          ts.unshift(t);
        }
      });
    }
    this.tribes = ts;
    return ts;
  }

  bots: Bot[] = [];
  myBots: Bot[] = [];

  @action async getBots(uniqueName?: string, queryParams?: any): Promise<any> {
    console.log("get bots");

    let query = this.appendQueryParams("bots", queryParams)
    let b = await api.get(query);

    const info = uiStore.meInfo;

    if (uniqueName) {
      b.forEach(function (t: Bot, i: number) {
        if (t.unique_name === uniqueName) {
          b.splice(i, 1);
          b.unshift(t);
        }
      });
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

    const hideBots = ['pleaseprovidedocumentation', 'example']

    // hide test bots and set images
    b &&
      b.forEach((bb, i) => {
        if (bb.unique_name === "btc") {
          // bb.img = "/static/bots_bitcoin.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name === "bet") {
          // bb.img = "/static/bots_betting.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name === "hello" || bb.unique_name === "welcome") {
          // bb.img = "/static/bots_welcome.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name && (bb.unique_name.includes("test") || hideBots.includes(bb.unique_name))) {
          // hide all test bots
          bb.hide = true;
        }

        if (bb.owner_pubkey === info?.owner_pubkey) {
          // hide my own bots
          bb.hide = true;
        }
      });

    this.bots = b;
    return b;

  }

  @action async getMyBots(): Promise<any> {
    if (!uiStore.meInfo) return null;

    const info = uiStore.meInfo;
    try {
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
      let relayB: any = await fetch(URL + "/bots", {
        method: "GET",
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });

      relayB = await relayB.json()
      console.log("got bots from relay", relayB);
      let relayMyBots = relayB?.response?.bots || []

      // merge tribe server stuff
      console.log("get bots");
      let tribeServerBots = await api.get(`bots/owner/${info.owner_pubkey}`);

      // merge data from tribe server, it has more than relay
      let mergedBots = relayMyBots.map(b => {
        const thisBot = tribeServerBots.find(f => f.uuid === b.uuid)
        return {
          ...b,
          ...thisBot
        }
      })

      this.myBots = mergedBots

      return mergedBots
    } catch (e) {
      console.log('ok')
    }
  }




  @action async getTribesByOwner(pubkey: string): Promise<Tribe[]> {
    const ts = await api.get(`tribes_by_owner/${pubkey}?all=true`);
    this.ownerTribes = ts;
    return ts;
  }

  @action async getGithubIssueData(
    owner: string,
    repo: string,
    issue: string
  ): Promise<any> {
    const data = await api.get(`github_issue/${owner}/${repo}/${issue}`);
    const { title, description, assignee, status } = data && data;

    console.log('got github issue', data)

    // if no title, the github issue isnt real
    if (!title && !status && !description && !assignee) return null;
    return data;
  }

  @action async getOpenGithubIssues(): Promise<any> {
    try {
      const openIssues = await api.get(`github_issue/status/open`);
      console.log('got openIssues', openIssues)
      // remove my own!

      let oIssues = [...openIssues]
      if (openIssues) {
        // remove my issues from count
        // if (uiStore.meInfo?.owner_pubkey) oIssues.filter(f => f.owner_pubkey != uiStore.meInfo?.owner_pubkey)
        uiStore.setOpenGithubIssues(oIssues)
      }
      return oIssues;
    } catch (e) {
      console.log('e', e)
    }
  }

  @action async makeBot(payload: any): Promise<any> {

    if (!uiStore.meInfo) return null;
    const info = uiStore.meInfo;
    try {
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
      let b: any = await fetch(URL + "/bot", {
        method: "POST",
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...payload
        }),
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });


      b = await b.json()
      console.log("made bot", b);

      const mybots = await this.getMyBots()
      console.log("got my bots", mybots);

      return b?.response;
    } catch (e) {
      console.log('failed', e)
    }

  }

  @action async updateBot(payload: any): Promise<any> {

    if (!uiStore.meInfo) return null;
    const info = uiStore.meInfo;
    try {
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
      const b = await fetch(URL + "/bot", {
        method: "PUT",
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...payload
        }),
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });
      console.log("updated bot", b);

      return b;
    } catch (e) {
      console.log('ok')
    }

  }

  @action async deleteBot(id: string): Promise<any> {

    if (!uiStore.meInfo) return null;
    const info = uiStore.meInfo;

    // delete from relay
    try {
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
      const b = await fetch(URL + `/bot/${id}`, {
        method: "DELETE",
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
        }),
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });

      console.log("deleted from relay", b);

      return b;
    } catch (e) {
      console.log('failed!')
    }

  }

  @action async postToCache(payload: any): Promise<void> {
    await api.post("save", payload, {
      "Content-Type": "application/json",
    });
    return;
  }

  @action appendQueryParams(path: string, limit: number, queryParams?: QueryParams): string {
    let query = path
    if (queryParams) {
      queryParams.limit = limit
      query += '?'
      const length = Object.keys(queryParams).length
      Object.keys(queryParams).forEach((k, i) => {
        query += `${k}=${queryParams[k]}`

        // add & if not last param
        if (i !== length - 1) {
          query += '&'
        }
      })
    }

    return query
  }

  // @persist("list")
  @observable
  people: Person[] = [];

  @action async getPeople(queryParams?: any): Promise<Person[]> {

    queryParams = { ...queryParams, search: uiStore.searchText }

    let query = this.appendQueryParams("people", queryLimit, { ...queryParams, sortBy: 'updated' })
    let ps = await api.get(query);

    if (uiStore.meInfo) {
      const index = ps.findIndex((f) => f.id == uiStore.meInfo?.id);
      if (index > -1) {
        // add 'hide' property to me in people list
        ps[index].hide = true;
      }
    }

    // this is for ordering, fix me, search is its own query
    // if (uniqueName) {
    //   ps?.forEach(function (t: Tribe, i: number) {
    //     if (t.unique_name === uniqueName) {
    //       ps.splice(i, 1);
    //       ps.unshift(t);
    //     }
    //   });
    // }

    // console.log('ps', ps)

    // for search always reset page
    if (queryParams && queryParams.resetPage) {
      this.people = ps
      uiStore.setPeoplePageNumber(1)
    } else {
      // all other cases, merge
      this.people = this.doPageListMerger(
        this.people,
        ps,
        uiStore.peoplePageNumber,
        (n) => uiStore.setPeoplePageNumber(n),
        queryLimit,
        queryParams
      )
    }

    return ps;
  }

  @action decodeListJSON(li: any): Promise<any[]> {
    if (li?.length) {
      li.forEach((o, i) => {
        li[i].body = JSON.parse(o.body)
        li[i].person = JSON.parse(o.person)
      })
    }
    return li
  }

  // @persist("list")
  @observable
  peoplePosts: PersonPost[] = [];

  @action async getPeoplePosts(queryParams?: any): Promise<PersonPost[]> {
    // console.log('queryParams', queryParams)
    queryParams = { ...queryParams, search: uiStore.searchText }
    let query = this.appendQueryParams("people/posts", smallQueryLimit, { ...queryParams, sortBy: 'created' })
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps)

      // console.log('ps', ps)

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peoplePosts = ps
        uiStore.setPeoplePostsPageNumber(1)
      } else {
        // all other cases, merge
        this.peoplePosts = this.doPageListMerger(
          this.peoplePosts,
          ps,
          uiStore.peoplePostsPageNumber,
          (n) => uiStore.setPeoplePostsPageNumber(n),
          smallQueryLimit,
          queryParams)
      }
      return ps;
    } catch (e) {
      console.log('fetch failed', e)
      return [];
    }
  }

  // @persist("list")
  @observable
  peopleWanteds: PersonWanted[] = [];

  @action async getPeopleWanteds(queryParams?: any): Promise<PersonWanted[]> {
    // console.log('queryParams', queryParams)
    queryParams = { ...queryParams, search: uiStore.searchText }
    let query = this.appendQueryParams("people/wanteds", smallQueryLimit, { ...queryParams, sortBy: 'created' })
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps)

      // console.log('ps', ps)

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peopleWanteds = ps
        uiStore.setPeopleWantedsPageNumber(1)
      } else {
        // all other cases, merge
        this.peopleWanteds = this.doPageListMerger(
          this.peopleWanteds,
          ps,
          uiStore.peopleWantedsPageNumber,
          (n) => uiStore.setPeopleWantedsPageNumber(n),
          smallQueryLimit,
          queryParams)
      }
      return ps;
    } catch (e) {
      console.log('fetch failed', e)
      return [];
    }
  }

  // @persist("list")
  @observable
  peopleOffers: PersonOffer[] = [];

  @action async getPeopleOffers(queryParams?: any): Promise<PersonOffer[]> {
    // console.log('queryParams', queryParams)
    queryParams = { ...queryParams, search: uiStore.searchText }
    let query = this.appendQueryParams("people/offers", smallQueryLimit, { ...queryParams, sortBy: 'created' })
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps)

      // console.log('ps', ps)

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peopleOffers = ps
        uiStore.setPeopleOffersPageNumber(1)
      } else {
        // all other cases, merge
        this.peopleOffers = this.doPageListMerger(
          this.peopleOffers,
          ps,
          uiStore.peopleOffersPageNumber,
          (n) => uiStore.setPeopleOffersPageNumber(n),
          smallQueryLimit,
          queryParams)
      }

      return ps;
    } catch (e) {
      console.log('fetch failed', e)
      return [];
    }

  }


  @action doPageListMerger(currentList: any[], newList: any[], pageNumber: number, setPage: Function, limit: number, queryParams?: any) {
    if (!newList || !newList.length) {
      console.log('got no new items, do not change page')
      return currentList
    }

    // FIX ME, make me an infinite loader

    // console.log('newList', newList)

    // let whileIndex = 0
    // while (newList.length < limit) {
    //   newList.push({ created: (163463 * whileIndex + 1), hide: true, lastPage: true })
    //   whileIndex++
    // }

    let direction = 0

    if (queryParams?.page) {
      // this tells us whether we're loading earlier data, or later data so we can merge the array in the right order
      if (pageNumber < queryParams.page) direction = 1 // paging forward
      else direction = -1 // paging backward

      // update page number in ui
      setPage(queryParams.page)
    }

    let keepGroup = [...currentList]
    let merger
    // no page movement, all incoming are the new list
    // if (direction === 0) {
    //   merger = newList
    // }
    // // paging forward
    // else if (direction > 0) {
    //   if (keepGroup.length === limit * 2) {
    //     keepGroup = keepGroup.slice(limit, limit + limit)
    //   }
    //   merger = [...keepGroup, ...newList];
    // }
    // // paging backward
    // else {

    //   // check if last page
    //   // if (keepGroup.includes(f => f.lastPage)) {
    //   //   merger = [...newList];
    //   // } else {
    //   // keepGroup = keepGroup.slice(0, limit)
    //   merger = newList//[...newList, ...keepGroup];
    //   // }
    //   console.log('merger', merger)
    // }

    merger = newList

    // if (merger.length > limit * 2) {
    //   merger = merger.slice(0, limit * 2)
    // }

    // let ids: any = []
    // merger.forEach((p: any, i: number) => {
    //   if (!ids.includes(p.created)) ids.push(p.created)
    //   else {
    //     merger.splice(1, i)
    //   }
    // })

    return merger
  }

  @action async getPersonByPubkey(pubkey: string): Promise<Person> {
    const p = await api.get(`person/${pubkey}`);
    // console.log('p', p)
    return p
  }


  // this method merges the relay self data with the db self data, they each hold different data
  @action async getSelf(me: any) {
    console.log('getSelf')
    let self = me || uiStore.meInfo
    if (self) {
      const p = await api.get(`person/${self.owner_pubkey}`);

      let updateSelf = { ...self, ...p }
      console.log('updateSelf', updateSelf)
      uiStore.setMeInfo(updateSelf);
    }
  }

  @action async refreshJwt() {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
      const res: any = await fetch(URL + "/refresh_jwt", {
        method: "GET",
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });
      const j = await res.json();

      return j.response;
    } catch (e) {
      console.log("e", e);
      // could not refresh jwt, logout!
      return null;
    }
  }

  @action async getUsdToSatsExchangeRate() {
    try {
      // get rate for 1 USD
      const res: any = await fetch(
        "https://blockchain.info/tobtc?currency=USD&value=1",
        {
          method: "GET",
        }
      );
      const j = await res.json();
      // 1 bitcoin is 1 million satoshis
      let satoshisInABitcoin = 0.00000001;
      const exchangeRate = j / satoshisInABitcoin;

      console.log("update exchange rate", exchangeRate);
      uiStore.setUsdToSatsExchangeRate(exchangeRate);

      return exchangeRate;
    } catch (e) {
      console.log("e", e);
      // could not refresh jwt, logout!
      return null;
    }
  }

  @action async deleteProfile() {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
      const res: any = await fetch(URL + "/profile", {
        method: "DELETE",
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...info,
        }),
      });

      uiStore.setMeInfo(null);
      uiStore.setSelectingPerson(0);
      uiStore.setSelectedPerson(0);

      const j = await res.json();
      return j;
    } catch (e) {
      console.log("e", e);
      // could not delete profile!
      return null;
    }
  }

  @action async addFavorite() {
    let body: any = {};
    console.log("SUBMIT FORM", body);

    // console.log('mergeFormWithMeData', body);
    if (!body) return; // avoid saving bad state

    const info = uiStore.meInfo as any;
    if (!info) return console.log("no meInfo");
    try {
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
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
      });

      if (!r.ok) {
        return alert("Failed to save data");
      }

      uiStore.setToasts([
        {
          id: "1",
          title: "Added to favorites.",
        },
      ]);
    } catch (e) {
      console.log("e", e);
    }
  }

  @action async deleteFavorite() {
    let body: any = {};
    console.log("SUBMIT FORM", body);

    // console.log('mergeFormWithMeData', body);
    if (!body) return; // avoid saving bad state

    const info = uiStore.meInfo as any;
    if (!info) return console.log("no meInfo");
    try {
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
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
      });

      if (!r.ok) {
        return alert("Failed to save data");
      }

      uiStore.setToasts([
        {
          id: "1",
          title: "Added to favorites.",
        },
      ]);
    } catch (e) {
      console.log("e", e);
    }
  }
}

export const mainStore = new MainStore();

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
  pubkey: string;
  photo_url: string;
  alias: string;
  route_hint: string;
  contact_key: string;
  price_to_meet: number;
  url: string;
  verification_signature: string;
  extras: Extras;
  hide?: boolean;
}

export interface PersonFlex {
  id?: number;
  unique_name?: string;
  owner_pubkey?: string;
  owner_alias?: string;
  description?: string;
  img?: string;
  tags?: string[];
  pubkey?: string;
  photo_url?: string;
  alias?: string;
  route_hint?: string;
  contact_key?: string;
  price_to_meet?: number;
  url?: string;
  verification_signature?: string;
  extras?: Extras;
  hide?: boolean;
}

export interface PersonPost {
  person: PersonFlex;
  title?: string;
  description?: string;
  created: number;
}

export interface PersonWanted {
  person: PersonFlex;
  title?: string;
  description?: string;
  created: number;
}

export interface PersonOffer {
  person: PersonFlex;
  title: string;
  description: string;
  created: number;
}

export interface Jwt {
  jwt: string;
}

export interface QueryParams {
  page?: number;
  limit?: number;
  sortBy?: string;
  direction?: string;
}
