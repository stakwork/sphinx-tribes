import { observable, action } from 'mobx';
// import { persist } from "mobx-persist";
import api from '../api';
import { Extras } from '../form/inputs/widgets/interfaces';
import { getHostIncludingDockerHosts } from '../host';
import { uiStore } from './ui';
import { randomString } from '../helpers';
export const queryLimit = 100;

function makeTorSaveURL(host: string, key: string) {
  return `sphinx.chat://?action=save&host=${host}&key=${key}`;
}

export class MainStore {
  // @persist("list")
  @observable
  tribes: Tribe[] = [];
  ownerTribes: Tribe[] = [];

  @action async getTribes(queryParams?: any): Promise<Tribe[]> {
    let ta = [...uiStore.tags];

    console.log('getTribes');
    //make tags string for querys
    ta = ta.filter((f) => f.checked);
    let tags = '';
    if (ta && ta.length) {
      ta.forEach((o, i) => {
        tags += o.label;
        if (ta.length - 1 !== i) {
          tags += ',';
        }
      });
    }
    queryParams = { ...queryParams, search: uiStore.searchText, tags };

    const query = this.appendQueryParams('tribes', queryLimit, {
      ...queryParams,
      sortBy: 'last_active=0, last_active',
      direction: 'desc'
    });
    const ts = await api.get(query);

    this.tribes = this.doPageListMerger(
      this.tribes,
      ts,
      (n) => uiStore.setTribesPageNumber(n),
      queryParams
    );

    return ts;
  }

  bots: Bot[] = [];
  myBots: Bot[] = [];

  @action async getBots(uniqueName?: string, queryParams?: any): Promise<any> {
    console.log('get bots');

    const query = this.appendQueryParams('bots', queryParams);
    const b = await api.get(query);

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

    const hideBots = ['pleaseprovidedocumentation', 'example'];

    // hide test bots and set images
    b &&
      b.forEach((bb, i) => {
        if (bb.unique_name === 'btc') {
          // bb.img = "/static/bots_bitcoin.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name === 'bet') {
          // bb.img = "/static/bots_betting.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (bb.unique_name === 'hello' || bb.unique_name === 'welcome') {
          // bb.img = "/static/bots_welcome.png";
          b.splice(i, 1);
          b.unshift(bb);
        }
        if (
          bb.unique_name &&
          (bb.unique_name.includes('test') || hideBots.includes(bb.unique_name))
        ) {
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
      let relayB: any = await this.fetchFromRelay('bots');

      relayB = await relayB.json();
      console.log('got bots from relay', relayB);
      const relayMyBots = relayB?.response?.bots || [];

      // merge tribe server stuff
      console.log('get bots');
      const tribeServerBots = await api.get(`bots/owner/${info.owner_pubkey}`);

      // merge data from tribe server, it has more than relay
      const mergedBots = relayMyBots.map((b) => {
        const thisBot = tribeServerBots.find((f) => f.uuid === b.uuid);
        return {
          ...b,
          ...thisBot
        };
      });

      this.myBots = mergedBots;

      return mergedBots;
    } catch (e) {
      console.log('ok');
    }
  }

  @action async fetchFromRelay(path): Promise<any> {
    if (!uiStore.meInfo) return null;

    const info = uiStore.meInfo;
    const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;
    const r: any = await fetch(`${URL}/${path}`, {
      method: 'GET',
      headers: {
        'x-jwt': info.jwt,
        'Content-Type': 'application/json'
      }
    });

    return r;
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
    const tribesClone = [...this.tribes];
    const dupIndex = tribesClone.findIndex((f) => f.uuid === t.uuid);
    if (dupIndex > -1) {
      tribesClone.splice(dupIndex, 1);
    }

    this.tribes = [t, ...tribesClone];
    return t;
  }

  @action async getSingleTribeByUn(un: string): Promise<Tribe> {
    const t = await api.get(`tribe_by_un/${un}`);
    return t;
  }

  @action async getGithubIssueData(owner: string, repo: string, issue: string): Promise<any> {
    const data = await api.get(`github_issue/${owner}/${repo}/${issue}`);
    const { title, description, assignee, status } = data && data;

    console.log('got github issue', data);

    // if no title, the github issue isnt real
    if (!title && !status && !description && !assignee) return null;
    return data;
  }

  @action async getOpenGithubIssues(): Promise<any> {
    try {
      const openIssues = await api.get(`github_issue/status/open`);
      console.log('got openIssues', openIssues);
      if (openIssues) {
        uiStore.setOpenGithubIssues(openIssues);
      }
      return openIssues;
    } catch (e) {
      console.log('e', e);
    }
  }

  @action isTorSave() {
    let result = false;
    if (uiStore?.meInfo?.url?.includes('.onion')) result = true;
    return result;
  }

  @action async makeBot(payload: any): Promise<any> {
    const [r, error] = await this.doCallToRelay('POST', `bot`, payload);
    if (error) throw error;
    if (!r) return; // tor user will return here

    const b = await r.json();
    console.log('made bot', b);

    const mybots = await this.getMyBots();
    console.log('got my bots', mybots);

    return b?.response;
  }

  @action async updateBot(payload: any): Promise<any> {
    const [r, error] = await this.doCallToRelay('PUT', `bot`, payload);
    if (error) throw error;
    if (!r) return; // tor user will return here
    console.log('updated bot', r);
    return r;
  }

  @action async deleteBot(id: string): Promise<any> {
    try {
      const [r, error] = await this.doCallToRelay('DELETE', `bot/${id}`, null);
      if (error) throw error;
      if (!r) return; // tor user will return here
      console.log('deleted from relay', r);
      return r;
    } catch (e) {
      console.log('failed!');
    }
  }

  @action async awardBadge(
    userPubkey: string,
    badgeName: string,
    badgeIcon: string,
    memo: string,
    amount?: number
  ): Promise<any> {
    const URL = 'https://liquid.sphinx.chat';
    let error;

    const info = uiStore.meInfo as any;
    if (!info) {
      error = new Error('Youre not logged in');
      return [null, error];
    }

    const headers = {
      'x-jwt': info.jwt,
      'Content-Type': 'application/json'
    };

    try {
      // 1. get user liquid address
      const userLiquidAddress = await api.get(`liquidAddressByPubkey/${userPubkey}`);

      if (!userLiquidAddress) {
        throw new Error('No Liquid Address tied to user account');
      }

      // 2. get password for login, login to "token" aliased as "tt"
      const res0 = await fetch(`${URL}/login`, {
        method: 'POST',
        body: JSON.stringify({
          pwd: 'password i got from user'
        }),
        headers
      });

      const j = await res0.json();
      const tt = j.token || '';

      // 3. first create the badge
      const res1 = await fetch(`${URL}/issue?token=${tt}`, {
        method: 'POST',
        body: JSON.stringify({
          name: badgeName,
          icon: badgeIcon,
          amount: amount || 1
        }),
        headers
      });

      const createdBadge = await res1.json();

      // 4. then transfer it
      const res2 = await fetch(`${URL}/transfer?token=${tt}`, {
        method: 'POST',
        body: JSON.stringify({
          asset: createdBadge.id,
          to: userLiquidAddress,
          amount: amount || 1,
          memo: memo || '1'
        }),
        headers
      });

      const transferredBadge = await res2.json();

      return transferredBadge;
    } catch (e) {
      console.log('ok');
    }
  }

  @action async getBadgeList(): Promise<any> {
    try {
      const URL = 'https://liquid.sphinx.chat';

      const l = await fetch(`${URL}/list`, {
        method: 'GET'
      });

      const badgelist = await l.json();

      uiStore.setBadgeList(badgelist);
      return badgelist;
    } catch (e) {
      console.log('ok');
    }
  }

  @action async getBalances(pubkey: any): Promise<any> {
    try {
      const URL = 'https://liquid.sphinx.chat';

      const b = await fetch(`${URL}/balances?pubkey=${pubkey}`, {
        method: 'GET'
      });

      const balances = await b.json();

      return balances;
    } catch (e) {
      console.log('ok');
    }
  }

  @action async postToCache(payload: any): Promise<void> {
    await api.post('save', payload, {
      'Content-Type': 'application/json'
    });
    return;
  }

  @action async getTorSaveURL(method: string, path: string, body: any): Promise<string> {
    const key = randomString(15);
    const gotHost = getHostIncludingDockerHosts();

    // make price to meet an integer
    if (body.price_to_meet) body.price_to_meet = parseInt(body.price_to_meet);

    const data = JSON.stringify({
      host: gotHost,
      ...body
    });

    let torSaveURL = '';

    try {
      await this.postToCache({
        key,
        body: data,
        path,
        method
      });
      torSaveURL = makeTorSaveURL(gotHost, key);
    } catch (e) {
      console.log('e', e);
    }

    return torSaveURL;
  }

  @action appendQueryParams(path: string, limit: number, queryParams?: QueryParams): string {
    let query = path;
    if (queryParams) {
      queryParams.limit = limit;
      query += '?';
      const { length } = Object.keys(queryParams);
      Object.keys(queryParams).forEach((k, i) => {
        query += `${k}=${queryParams[k]}`;

        // add & if not last param
        if (i !== length - 1) {
          query += '&';
        }
      });
    }

    return query;
  }

  @action async getPeopleByNameAliasPubkey(alias: string): Promise<Person[]> {
    const smallQueryLimit = 4;
    const query = this.appendQueryParams('people/search', smallQueryLimit, {
      search: alias,
      sortBy: 'owner_alias'
    });
    const ps = await api.get(query);
    return ps;
  }

  // @persist("list")
  @observable
  people: Person[] = [];

  @action async getPeople(queryParams?: any): Promise<Person[]> {
    queryParams = { ...queryParams, search: uiStore.searchText };

    const query = this.appendQueryParams('people', queryLimit, {
      ...queryParams,
      sortBy: 'last_login'
    });

    const ps = await api.get(query);

    if (uiStore.meInfo) {
      const index = ps.findIndex((f) => f.id === uiStore.meInfo?.id);
      if (index > -1) {
        // add 'hide' property to me in people list
        ps[index].hide = true;
      }
    }

    // for search always reset page
    if (queryParams && queryParams.resetPage) {
      this.people = ps;
      uiStore.setPeoplePageNumber(1);
    } else {
      // all other cases, merge
      this.people = this.doPageListMerger(
        this.people,
        ps,
        (n) => uiStore.setPeoplePageNumber(n),
        queryParams
      );
    }

    return ps;
  }

  @action decodeListJSON(li: any): Promise<any[]> {
    if (li?.length) {
      li.forEach((o, i) => {
        li[i].body = JSON.parse(o.body);
        li[i].person = JSON.parse(o.person);
      });
    }
    return li;
  }

  // @persist("list")
  @observable
  peoplePosts: PersonPost[] = [];

  @action async getPeoplePosts(queryParams?: any): Promise<PersonPost[]> {
    // console.log('queryParams', queryParams)
    queryParams = { ...queryParams, search: uiStore.searchText };

    const query = this.appendQueryParams('people/posts', queryLimit, {
      ...queryParams,
      sortBy: 'created'
    });
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps);

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peoplePosts = ps;
        uiStore.setPeoplePostsPageNumber(1);
      } else {
        // all other cases, merge
        this.peoplePosts = this.doPageListMerger(
          this.peoplePosts,
          ps,
          (n) => uiStore.setPeoplePostsPageNumber(n),
          queryParams
        );
      }
      return ps;
    } catch (e) {
      console.log('fetch failed', e);
      return [];
    }
  }

  // @persist("list")
  @observable
  peopleWanteds: PersonWanted[] = [];

  @action setPeopleWanteds(wanteds: PersonWanted[]) {
    this.peopleWanteds = wanteds;
  }

  @action async getPeopleWanteds(queryParams?: any): Promise<PersonWanted[]> {
    queryParams = { ...queryParams, search: uiStore.searchText };

    const query = this.appendQueryParams('people/wanteds', queryLimit, {
      ...queryParams,
      sortBy: 'created'
    });
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps);

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peopleWanteds = ps;
        uiStore.setPeopleWantedsPageNumber(1);
      } else {
        // all other cases, merge
        this.peopleWanteds = this.doPageListMerger(
          this.peopleWanteds,
          ps,
          (n) => uiStore.setPeopleWantedsPageNumber(n),
          queryParams
        );
      }
      return ps;
    } catch (e) {
      console.log('fetch failed', e);
      return [];
    }
  }

  // @persist("list")
  @observable
  peopleOffers: PersonOffer[] = [];

  @action async getPeopleOffers(queryParams?: any): Promise<PersonOffer[]> {
    // console.log('queryParams', queryParams)
    queryParams = { ...queryParams, search: uiStore.searchText };

    const query = this.appendQueryParams('people/offers', queryLimit, {
      ...queryParams,
      sortBy: 'created'
    });
    try {
      let ps = await api.get(query);
      ps = this.decodeListJSON(ps);

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.peopleOffers = ps;
        uiStore.setPeopleOffersPageNumber(1);
      } else {
        // all other cases, merge
        this.peopleOffers = this.doPageListMerger(
          this.peopleOffers,
          ps,
          (n) => uiStore.setPeopleOffersPageNumber(n),
          queryParams
        );
      }

      return ps;
    } catch (e) {
      console.log('fetch failed', e);
      return [];
    }
  }

  @action doPageListMerger(
    currentList: any[],
    newList: any[],
    setPage: Function,
    queryParams?: any
  ) {
    if (!newList || !newList.length) {
      if (queryParams.search) {
        // if search and no results, return nothing
        return [];
      } else {
        return currentList;
      }
    }

    if (queryParams && queryParams.resetPage) {
      setPage(1);
      return newList;
    }

    if (queryParams?.page) setPage(queryParams.page);
    const l = [...currentList, ...newList];

    return l;
  }

  @action async getPersonByPubkey(pubkey: string): Promise<Person> {
    const p = await api.get(`person/${pubkey}`);
    // console.log('p', p)
    return p;
  }

  @action async getPersonByGithubName(github: string): Promise<Person> {
    const p = await api.get(`person/githubname/${github}`);
    console.log('getPersonByGithubName', p);
    return p;
  }

  // this method merges the relay self data with the db self data, they each hold different data

  @action async getSelf(me: any) {
    console.log('getSelf');
    const self = me || uiStore.meInfo;
    if (self) {
      const p = await api.get(`person/${self.owner_pubkey}`);

      // get request for super_admin_array.
      const getSuperAdmin = async () => {
        try {
          const response = await api.get(`admin_pubkeys`);
          const admin_keys = response?.pubkeys;
          if (admin_keys !== null) {
            return !!admin_keys.find((value) => value === self.owner_pubkey);
          } else {
            return false;
          }
        } catch (error) {
          return false;
        }
      };

      const isSuperAdmin = await getSuperAdmin();
      const updateSelf = { ...self, ...p, isSuperAdmin: isSuperAdmin };
      console.log('updateSelf', updateSelf);
      uiStore.setMeInfo(updateSelf);
    }
  }

  @action async claimBadgeOnLiquid(body: ClaimOnLiquid): Promise<any> {
    try {
      const [r, error] = await this.doCallToRelay('POST', 'claim_on_liquid', body);
      if (error) throw error;
      if (!r) return; // tor user will return here

      console.log('code from relay', r);
      return r;
    } catch (e) {
      console.log('failed!', e);
    }
  }

  @action async sendBadgeOnLiquid(body: ClaimOnLiquid): Promise<any> {
    try {
      const [r, error] = await this.doCallToRelay('POST', 'claim_on_liquid', body);
      if (error) throw error;
      if (!r) return; // tor user will return here

      console.log('code from relay', r);
      return r;
    } catch (e) {
      console.log('failed!', e);
    }
  }

  @action async refreshJwt() {
    try {
      if (!uiStore.meInfo) return null;

      const res: any = await this.fetchFromRelay('refresh_jwt');
      const j = await res.json();

      return j.response;
    } catch (e) {
      console.log('e', e);
      // could not refresh jwt, logout!
      return null;
    }
  }

  @action async getUsdToSatsExchangeRate() {
    try {
      // get rate for 1 USD
      const res: any = await fetch('https://blockchain.info/tobtc?currency=USD&value=1', {
        method: 'GET'
      });
      const j = await res.json();
      // 1 bitcoin is 1 million satoshis
      const satoshisInABitcoin = 0.00000001;
      const exchangeRate = j / satoshisInABitcoin;

      console.log('update exchange rate', exchangeRate);
      uiStore.setUsdToSatsExchangeRate(exchangeRate);

      return exchangeRate;
    } catch (e) {
      console.log('e', e);
      // could not refresh jwt, logout!
      return null;
    }
  }

  @action async deleteProfile() {
    try {
      const info = uiStore.meInfo;
      const [r, error] = await this.doCallToRelay('DELETE', 'profile', info);
      if (error) throw error;
      if (!r) return; // tor user will return here

      uiStore.setMeInfo(null);
      uiStore.setSelectingPerson(0);
      uiStore.setSelectedPerson(0);

      const j = await r.json();
      return j;
    } catch (e) {
      console.log('e', e);
      // could not delete profile!
      return null;
    }
  }

  @action async saveProfile(body) {
    console.log('SUBMIT FORM', body);
    if (!body) return; // avoid saving bad state
    if (body.price_to_meet) body.price_to_meet = parseInt(body.price_to_meet); // must be an int

    try {
      const [r, error] = await this.doCallToRelay('POST', 'profile', body);
      if (error) throw error;
      if (!r) return; // tor user will return here

      // first time profile makers will need this on first login
      if (!body.id) {
        const j = await r.json();
        if (j.response.id) {
          body.id = j.response.id;
        }
      }

      uiStore.setToasts([
        {
          id: '1',
          title: 'Saved.'
        }
      ]);

      await this.getSelf(body);
    } catch (e) {
      console.log('e', e);
    }
  }

  // this method is used whenever changing data from the frontend,
  // forks between tor users and non-tor
  @action async doCallToRelay(method: string, path: string, body: any): Promise<any> {
    let response: Response;
    let error: any = null;

    const info = uiStore.meInfo as any;
    if (!info) {
      error = new Error('Youre not logged in');
      return [null, error];
    }

    // fork between tor users and not
    if (this.isTorSave()) {
      this.submitFormViaApp(method, path, body);
      return [null, null];
    }

    const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;

    response = await fetch(`${URL}/${path}`, {
      method: method,
      body: JSON.stringify({
        // use docker host (tribes.sphinx), because relay will post to it
        host: getHostIncludingDockerHosts(),
        ...body
      }),
      headers: {
        'x-jwt': info.jwt,
        'Content-Type': 'application/json'
      }
    });

    return [response, error];
  }

  @action async submitFormViaApp(method: string, path: string, body: any) {
    try {
      const torSaveURL = await this.getTorSaveURL(method, path, body);
      uiStore.setTorFormBodyQR(torSaveURL);
    } catch (e) {
      console.log('e', e);
    }
  }

  @action async setExtrasPropertyAndSave(
    extrasName: string,
    propertyName: string,
    created: number,
    newPropertyValue: any
  ): Promise<any> {
    if (uiStore.meInfo) {
      const clonedMeInfo = { ...uiStore.meInfo };
      const clonedExtras = clonedMeInfo?.extras;
      const clonedEx: any = clonedExtras && clonedExtras[extrasName];
      const targetIndex = clonedEx?.findIndex((f) => f.created === created);

      if (clonedEx && (targetIndex || targetIndex === 0) && targetIndex > -1) {
        try {
          clonedEx[targetIndex][propertyName] = newPropertyValue;
          clonedMeInfo.extras[extrasName] = clonedEx;
          await this.saveProfile(clonedMeInfo);
          return [clonedEx, targetIndex];
        } catch (e) {
          console.log('e', e);
        }
      }

      return [null, null];
    }
  }

  @action async deleteFavorite() {
    const body: any = {};
    console.log('SUBMIT FORM', body);

    // console.log('mergeFormWithMeData', body);
    if (!body) return; // avoid saving bad state

    const info = uiStore.meInfo as any;
    if (!info) return console.log('no meInfo');
    try {
      const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;
      const r = await fetch(`${URL}/profile`, {
        method: 'POST',
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...body,
          price_to_meet: parseInt(body.price_to_meet)
        }),
        headers: {
          'x-jwt': info.jwt,
          'Content-Type': 'application/json'
        }
      });

      if (!r.ok) {
        return alert('Failed to save data');
      }

      uiStore.setToasts([
        {
          id: '1',
          title: 'Added to favorites.'
        }
      ]);
    } catch (e) {
      console.log('e', e);
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
  show?: boolean;
  body: PersonWanted;
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
  search?: string;
}

export interface ClaimOnLiquid {
  asset: number;
  to: string;
  amount?: number;
  memo: string;
}
