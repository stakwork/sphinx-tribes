import { observable, action } from "mobx";
import { persist } from "mobx-persist";
import api from "../api";
import { Extras } from "../form/inputs/widgets/interfaces";
import { getHostIncludingDockerHosts } from "../host";
import { MeInfo, uiStore } from "./ui";

export const queryLimit = 5

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

    // hide test bots and set images
    b &&
      b.forEach((bb, i) => {
        if (bb.unique_name === "btc") {
          bb.img = "/static/bots_bitcoin.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name === "bet") {
          bb.img = "/static/bots_betting.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name === "hello" || bb.unique_name === "welcome") {
          bb.img = "/static/bots_welcome.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name && bb.unique_name.includes("test")) {
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
    const { title } = data && data;

    // if no title, the github issue isnt real
    if (!title) return null;
    return data;
  }

  @action async getOpenGithubIssues(): Promise<any> {
    try {
      const openIssues = await api.get(`github_issue/status/open`);
      return openIssues;
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

      return b?.response;
    } catch (e) {
      console.log('ok')
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
      console.log("deleted bot", b);

      return b;
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

  @action appendQueryParams(path: string, queryParams?: QueryParams): string {
    let query = path
    if (queryParams) {
      queryParams.limit = queryLimit
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

  @persist("list")
  @observable
  people: Person[] = [];

  @action async getPeople(uniqueName?: string, queryParams?: any): Promise<Person[]> {
    console.log('queryParams', queryParams)
    let query = this.appendQueryParams("people", { ...queryParams, sortBy: 'updated' })
    let ps = await api.get(query);

    if (!ps || !ps.length) {
      console.log('got no people, do not change page')
      return []
    }

    let direction = 0

    if (queryParams?.page) {
      // this tells us whether we're loading earlier data, or later data so we can merge the array in the right order
      if (uiStore.peoplePageNumber < queryParams.page) direction = 1 // paging forward
      else direction = -1 // paging backward

      // update people page in ui
      uiStore.setPeoplePageNumber(queryParams.page)
    }

    // fixme, this is old, dont need to do this if getSelf is updating properly
    if (uiStore.meInfo) {
      const index = ps.findIndex((f) => f.id == uiStore.meInfo?.id);

      if (index > -1) {
        // add 'hide' property to me in people list
        ps[index].hide = true;

        let meInfoUpdates: any = {};

        // if meInfo has no img but people list does, copy that
        if (!uiStore.meInfo.img && ps[index].img) {
          meInfoUpdates.img = ps[index].img;
        }

        // if meInfo has no github_issues but people list does, copy that
        if (!uiStore.meInfo.github_issues && ps[index].github_issues) {
          meInfoUpdates.github_issues = ps[index].github_issues;
        }

        // if meInfo has no verification_signature but people list does, copy that
        if (
          !uiStore.meInfo.verification_signature &&
          ps[index].verification_signature
        ) {
          meInfoUpdates.verification_signature =
            ps[index].verification_signature;
        }

        uiStore.setMeInfo({ ...uiStore.meInfo, ...meInfoUpdates });
      }
    }

    if (uniqueName) {
      ps.forEach(function (t: Tribe, i: number) {
        if (t.unique_name === uniqueName) {
          ps.splice(i, 1);
          ps.unshift(t);
        }
      });
    }

    let keepGroup = [...this.people]
    let mergePeople
    if (direction === 0) {
      mergePeople = ps
    }
    // paging forward
    else if (direction > 0) {
      keepGroup = keepGroup.slice(0, queryLimit - 1)
      mergePeople = [...keepGroup, ...ps];
    }
    // paging backward
    else {
      keepGroup = keepGroup.slice(queryLimit, queryLimit + queryLimit - 1)
      mergePeople = [...ps, ...keepGroup];
    }

    console.log('mergePeople', mergePeople)

    // remove duplicates if any
    let ids: any = []
    mergePeople.forEach((p: any, i: number) => {
      if (!ids.includes(p.id)) ids.push(p.id)
      else {
        console.log('found duplicates!', p.id)
        mergePeople[i].hide = true
      }
    })

    this.people = mergePeople

    return ps;
  }

  @action async getPersonByPubkey(pubkey: string): Promise<Person> {
    const p = await api.get(`person/${pubkey}`);
    console.log('p', p)
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

export interface Jwt {
  jwt: string;
}

export interface QueryParams {
  page?: number;
  limit?: number;
  sortBy?: string;
  direction?: string;
}
