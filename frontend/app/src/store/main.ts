import { observable, action } from "mobx";
import { persist } from "mobx-persist";
import api from "../api";
import { Extras } from "../form/inputs/widgets/interfaces";
import { getHostIncludingDockerHosts } from "../host";
import { uiStore } from "./ui";
import { randomString } from "../helpers";
export const queryLimit = 100

function makeTorSaveURL(host: string, key: string) {
  return `sphinx.chat://?action=save&host=${host}&key=${key}`;
}

export class MainStore {
  // @persist("list")
  @observable
  tribes: Tribe[] = [];
  ownerTribes: Tribe[] = [];

  @action async getTribes(queryParams?: any): Promise<Tribe[]> {
    let ta = [...uiStore.tags]

    console.log('getTribes')
    //make tags string for querys
    ta = ta.filter(f => f.checked)
    let tags = ''
    if (ta && ta.length) {
      ta.forEach((o, i) => {
        tags += o.label
        if (ta.length - 1 !== i) {
          tags += ','
        }
      })
    }
    queryParams = { ...queryParams, search: uiStore.searchText, tags }

    let query = this.appendQueryParams("tribes", queryLimit, { ...queryParams, sortBy: 'last_active=0, last_active', direction: 'desc' })
    const ts = await api.get(query);

    this.tribes = this.doPageListMerger(
      this.tribes,
      ts,
      (n) => uiStore.setTribesPageNumber(n),
      queryParams)

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

  @action async getTribeByUn(un: string): Promise<Tribe> {
    const t = await api.get(`tribe_by_un/${un}`);
    // put got on top
    // if already exists, delete
    const tribesClone = [...this.tribes]
    const dupIndex = tribesClone.findIndex(f => f.uuid === t.uuid)
    if (dupIndex > -1) {
      tribesClone.splice(dupIndex, 1)
    }

    this.tribes = [t, ...tribesClone]
    return t;
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
      // console.log('got openIssues', openIssues)
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

  @action async getBadgeList(): Promise<any> {

    try {
      const URL = 'https://liquid.sphinx.chat'

      const l = await fetch(URL + `/list`, {
        method: "GET",
      });

      const badgelist = await l.json();

      uiStore.setBadgeList(badgelist)
      return badgelist;
    } catch (e) {
      console.log('ok')
    }

  }

  @action async getBalances(pubkey: any): Promise<any> {

    try {
      const URL = 'https://liquid.sphinx.chat'

      const b = await fetch(URL + `/balances?pubkey=${pubkey}`, {
        method: "GET",
      });

      const balances = await b.json();

      return balances;
    } catch (e) {
      console.log('ok')
    }

  }

  @action async postToCache(payload: any): Promise<void> {
    await api.post("save", payload, {
      "Content-Type": "application/json",
    });
    return;
  }

  @action async getTorSaveURL(method: string, path: string, body: any): Promise<string> {
    const key = randomString(15);
    const gotHost = getHostIncludingDockerHosts();

    // make price to meet an integer
    if (body.price_to_meet) body.price_to_meet = parseInt(body.price_to_meet)

    const data = JSON.stringify({
      host: gotHost,
      ...body
    });

    let torSaveURL = ''

    try {
      await this.postToCache({
        key,
        body: data,
        path,
        method,
      });
      torSaveURL = makeTorSaveURL(gotHost, key);
    } catch (e) {
      console.log('e', e)
    }

    return torSaveURL
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

    let query = this.appendQueryParams("people", queryLimit, { ...queryParams, sortBy: 'last_login' })

    let ps = await api.get(query);

    if (uiStore.meInfo) {
      const index = ps.findIndex((f) => f.id == uiStore.meInfo?.id);
      if (index > -1) {
        // add 'hide' property to me in people list
        ps[index].hide = true;
      }
    }

    // for search always reset page
    if (queryParams && queryParams.resetPage) {
      this.people = ps
      uiStore.setPeoplePageNumber(1)
    } else {
      // all other cases, merge
      this.people = this.doPageListMerger(
        this.people,
        ps,
        (n) => uiStore.setPeoplePageNumber(n),
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

    let query = this.appendQueryParams("people/posts", queryLimit, { ...queryParams, sortBy: 'created' })
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps)

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peoplePosts = ps
        uiStore.setPeoplePostsPageNumber(1)
      } else {
        // all other cases, merge
        this.peoplePosts = this.doPageListMerger(
          this.peoplePosts,
          ps,
          (n) => uiStore.setPeoplePostsPageNumber(n),
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

  @action setPeopleWanteds(wanteds: PersonWanted[]) {
    this.peopleWanteds = wanteds
  }

  @action async getPeopleWanteds(queryParams?: any): Promise<PersonWanted[]> {
    queryParams = { ...queryParams, search: uiStore.searchText }

    let query = this.appendQueryParams("people/wanteds", queryLimit, { ...queryParams, sortBy: 'created' })
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps)

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peopleWanteds = ps
        uiStore.setPeopleWantedsPageNumber(1)
      } else {
        // all other cases, merge
        this.peopleWanteds = this.doPageListMerger(
          this.peopleWanteds,
          ps,
          (n) => uiStore.setPeopleWantedsPageNumber(n),
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

    let query = this.appendQueryParams("people/offers", queryLimit, { ...queryParams, sortBy: 'created' })
    try {

      let ps = await api.get(query);
      ps = this.decodeListJSON(ps)

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peopleOffers = ps
        uiStore.setPeopleOffersPageNumber(1)
      } else {
        // all other cases, merge
        this.peopleOffers = this.doPageListMerger(
          this.peopleOffers,
          ps,
          (n) => uiStore.setPeopleOffersPageNumber(n),
          queryParams)
      }

      return ps;
    } catch (e) {
      console.log('fetch failed', e)
      return [];
    }

  }

  @action doPageListMerger(currentList: any[], newList: any[], setPage: Function, queryParams?: any) {
    if (!newList || !newList.length) {
      if (queryParams.search) {
        // if search and no results, return nothing
        return []
      } else {
        return currentList
      }
    }

    if (queryParams && queryParams.resetPage) {
      setPage(1)
      return newList
    }

    if (queryParams?.page) setPage(queryParams.page)
    let l = [...currentList, ...newList]

    return l
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



  @action async claimBadgeOnLiquid(body: ClaimOnLiquid): Promise<any> {

    if (!uiStore.meInfo) return null;
    const info = uiStore.meInfo;

    try {
      const URL = info.url.startsWith("http")
        ? info.url
        : `https://${info.url}`;
      const b = await fetch(URL + `/claim_on_liquid`, {
        method: "POST",
        body: JSON.stringify({
          ...body,
          host: getHostIncludingDockerHosts(),
        }),
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });

      console.log("code from relay", b);

      return b;
    } catch (e) {
      console.log('failed!', e)
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

  @action async saveProfile(body) {
    console.log("SUBMIT FORM", body);

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
          title: "Saved.",
        },
      ]);

      // await this.getSelf(body);
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
  last_login?: number;
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
  last_login?: number;
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
  show: boolean;
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


export interface ClaimOnLiquid {
  asset: number
  to: string
  amount?: number
  memo: string
}