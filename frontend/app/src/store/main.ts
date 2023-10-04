import { makeAutoObservable, observable, action } from 'mobx';
import memo from 'memo-decorator';
import { persist } from 'mobx-persist';
import { uniqBy } from 'lodash';
import api from '../api';
import { Extras } from '../components/form/inputs/widgets/interfaces';
import { getHostIncludingDockerHosts } from '../config/host';
import { randomString } from '../helpers';
import { TribesURL } from '../config/host';
import { uiStore } from './ui';
import { getUserAvatarPlaceholder } from './lib';

export const queryLimit = 100;

function makeTorSaveURL(host: string, key: string) {
  return `sphinx.chat://?action=save&host=${host}&key=${key}`;
}

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
  pubkey?: string;
  photo_url: string;
  alias: string;
  route_hint: string;
  owner_route_hint?: string;
  contact_key: string;
  price_to_meet: number;
  last_login?: number;
  url: string;
  verification_signature: string;
  extras: Extras;
  hide?: boolean;
  commitment_fee?: number;
  assigned_hours?: number;
  bounty_expires?: number;
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
  person?: any;
  body?: any;
  org_uuid?: any;
  title?: string;
  description?: string;
  owner_id: string;
  created?: number;
  show?: boolean;
  assignee?: any;
  wanted_type: string;
  type?: string;
  price?: string;
  codingLanguage: string;
  estimated_session_length: string;
  bounty_expires?: string;
  commitment_fee?: number;
}

export interface PaymentHistory {
  id: number;
  bounty_id: number;
  amount: number;
  org_uuid: string;
  sender_name: string;
  receiver_name: string;
  created: string;
}

export interface BudgetHistory {
  id: number;
  amount: number;
  org_uuid: string;
  payment_type: string;
  created: string;
  updated: string;
  sender_pub_key: string;
  sender_name: string;
  status: boolean;
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
  resetPage?: boolean;
}

export interface ClaimOnLiquid {
  asset: number;
  to: string;
  amount?: number;
  memo: string;
}

export interface LnAuthData {
  encode: string;
  k1: string;
}

export interface LnInvoice {
  success: boolean;
  response: {
    invoice: string;
  };
}

export interface Organization {
  id: string;
  uuid: string;
  name: string;
  owner_pubkey: string;
  img: string;
  created: string;
  updated: string;
  show: boolean;
  bounty_count?: number;
}

export interface BountyRoles {
  name: string;
}
export class MainStore {
  [x: string]: any;
  tribes: Tribe[] = [];
  ownerTribes: Tribe[] = [];

  constructor() {
    makeAutoObservable(this);
  }

  async getTribes(queryParams?: any): Promise<Tribe[]> {
    let ta = [...uiStore.tags];

    //make tags string for querys
    ta = ta.filter((f: any) => f.checked);
    let tags = '';
    if (ta && ta.length) {
      ta.forEach((o: any, i: any) => {
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
      (n: any) => uiStore.setTribesPageNumber(n),
      queryParams
    );

    return ts;
  }

  bots: Bot[] = [];
  myBots: Bot[] = [];

  async getBots(uniqueName?: string, queryParams?: any): Promise<any> {
    const query = this.appendQueryParams('bots', 100, queryParams);
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

    const hideBots = ['pleaseprovidedocumentation', 'example'];

    // hide test bots and set images
    b &&
      b.forEach((bb: any, i: any) => {
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

  async getMyBots(): Promise<any> {
    if (!uiStore.meInfo) return null;

    const info = uiStore.meInfo;
    try {
      let relayB: any = await this.fetchFromRelay('bots');

      relayB = await relayB.json();
      const relayMyBots = relayB?.response?.bots || [];

      // merge tribe server stuff
      const tribeServerBots = await api.get(`bots/owner/${info.owner_pubkey}`);

      // merge data from tribe server, it has more than relay
      const mergedBots = relayMyBots.map((b: any) => {
        const thisBot = tribeServerBots.find((f: any) => f.uuid === b.uuid);
        return {
          ...b,
          ...thisBot
        };
      });

      this.myBots = mergedBots;

      return mergedBots;
    } catch (e) {
      console.log('Error getMyBots', e);
    }
  }

  async fetchFromRelay(path: string): Promise<any> {
    if (!uiStore.meInfo) return null;

    const info = uiStore.meInfo;
    const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;

    const r: any = await fetch(`${URL}/${path}`, {
      method: 'GET',
      mode: 'cors',
      headers: {
        'x-jwt': info.jwt,
        'Content-Type': 'application/json',
        Accept: 'application/json'
      }
    });

    return r;
  }

  async getTribesByOwner(pubkey: string): Promise<Tribe[]> {
    const ts = await api.get(`tribes_by_owner/${pubkey}?all=true`);
    this.ownerTribes = ts;
    return ts;
  }

  async getTribeByUn(un: string): Promise<Tribe> {
    const t = await api.get(`tribe_by_un/${un}`);
    // put got on top
    // if already exists, delete
    const tribesClone = [...this.tribes];
    const dupIndex = tribesClone.findIndex((f: any) => f.uuid === t.uuid);
    if (dupIndex > -1) {
      tribesClone.splice(dupIndex, 1);
    }

    this.tribes = [t, ...tribesClone];
    return t;
  }

  async getSingleTribeByUn(un: string): Promise<Tribe> {
    const t = await api.get(`tribe_by_un/${un}`);
    return t;
  }

  async getGithubIssueData(owner: string, repo: string, issue: string): Promise<any> {
    const data = await api.get(`github_issue/${owner}/${repo}/${issue}`);
    const { title, description, assignee, status } = data && data;

    // if no title, the github issue isnt real
    if (!title && !status && !description && !assignee) return null;
    return data;
  }

  async getOpenGithubIssues(): Promise<any> {
    try {
      const openIssues = await api.get(`github_issue/status/open`);
      if (openIssues) {
        uiStore.setOpenGithubIssues(openIssues);
      }
      return openIssues;
    } catch (e) {
      console.log('Error getOpenGithubIssues: ', e);
    }
  }

  isTorSave() {
    let result = false;
    if (uiStore?.meInfo?.url?.includes('.onion')) result = true;
    return result;
  }

  async makeBot(payload: any): Promise<any> {
    const [r, error] = await this.doCallToRelay('POST', `bot`, payload);
    if (error) throw error;
    if (!r) return; // tor user will return here

    const b = await r.json();

    // const mybots = await this.getMyBots();

    return b?.response;
  }

  async updateBot(payload: any): Promise<any> {
    const [r, error] = await this.doCallToRelay('PUT', `bot`, payload);
    if (error) throw error;
    if (!r) return; // tor user will return here
    return r;
  }

  async deleteBot(id: string): Promise<any> {
    try {
      const [r, error] = await this.doCallToRelay('DELETE', `bot/${id}`, null);
      if (error) throw error;
      if (!r) return; // tor user will return here
      return r;
    } catch (e) {
      console.log('Error deleteBot: ', e);
    }
  }

  async awardBadge(
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
      const tt = j.token || this.lnToken || '';

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
      console.log('Error awardBadge: ', e);
    }
  }

  async getBadgeList(): Promise<any> {
    try {
      const URL = 'https://liquid.sphinx.chat';

      const l = await fetch(`${URL}/list?limit=100000`, {
        method: 'GET'
      });

      const badgelist = await l.json();

      uiStore.setBadgeList(badgelist);
      return badgelist;
    } catch (e) {
      console.log('Error getBadgeList: ', e);
    }
  }

  async getBalances(pubkey: any): Promise<any> {
    try {
      const URL = 'https://liquid.sphinx.chat';

      const b = await fetch(`${URL}/balances?pubkey=${pubkey}&limit=100000`, {
        method: 'GET'
      });

      const balances = await b.json();

      return balances;
    } catch (e) {
      console.log('Error getBalances: ', e);
    }
  }

  async postToCache(payload: any): Promise<void> {
    await api.post('save', payload, {
      'Content-Type': 'application/json'
    });
    return;
  }

  async getTorSaveURL(method: string, path: string, body: any): Promise<string> {
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
      console.log('Error postToCache getTorSaveURL: ', e);
    }

    return torSaveURL;
  }

  appendQueryParams(path: string, limit: number, queryParams?: QueryParams): string {
    const adaptedParams = {
      ...queryParams,
      limit: String(limit),
      ...(queryParams?.resetPage ? { resetPage: String(queryParams.resetPage) } : {}),
      ...(queryParams?.page ? { page: String(queryParams.page) } : {})
    } as Record<string, string>;

    const searchParams = new URLSearchParams(adaptedParams);

    return `${path}?${searchParams.toString()}`;
  }

  async getPeopleByNameAliasPubkey(alias: string): Promise<Person[]> {
    const smallQueryLimit = 4;
    const query = this.appendQueryParams('people/search', smallQueryLimit, {
      search: alias,
      sortBy: 'owner_alias'
    });
    const ps = await api.get(query);
    return ps;
  }

  getUserAvatarPlaceholder(ownerId: string) {
    return getUserAvatarPlaceholder(ownerId);
  }

  @persist('list')
  _people: Person[] = [];

  get people() {
    return this._people.map((person: Person) => ({
      ...person,
      img: person.img || this.getUserAvatarPlaceholder(person.owner_pubkey)
    }));
  }

  set people(people: Person[]) {
    this._people = uniqBy(people, 'uuid');
  }

  setPeople(p: Person[]) {
    this._people = p;
  }

  async getPeople(queryParams?: any): Promise<Person[]> {
    const params = { ...queryParams, search: uiStore.searchText };
    const ps = await this.fetchPeople(uiStore.searchText, queryParams);

    if (uiStore.meInfo) {
      const index = ps.findIndex((f: any) => f.id === uiStore.meInfo?.id);
      if (index > -1) {
        // add 'hide' property to me in people list
        ps[index].hide = true;
      }
    }

    // for search always reset page
    if (params && params.resetPage) {
      this.people = ps;
      uiStore.setPeoplePageNumber(1);
    } else {
      // all other cases, merge
      this.people = this.doPageListMerger(
        this.people,
        ps,
        (n: any) => uiStore.setPeoplePageNumber(n),
        params
      );
    }

    return ps;
  }

  @memo({
    resolver: (...args: any[]) => JSON.stringify({ args }),
    cache: new Map()
  })
  private async fetchPeople(search: string, queryParams?: any): Promise<Person[]> {
    const params = { ...queryParams, search };
    const query = this.appendQueryParams('people', queryLimit, {
      ...params,
      sortBy: 'last_login'
    });
    const ps = await api.get(query);
    return ps;
  }

  decodeListJSON(li: any): Promise<any[]> {
    if (li?.length) {
      li.forEach((o: any, i: any) => {
        li[i].body = JSON.parse(o.body);
        li[i].person = JSON.parse(o.person);
      });
    }
    return li;
  }

  @persist('list')
  peoplePosts: PersonPost[] = [];

  async getPeoplePosts(queryParams?: any): Promise<PersonPost[]> {
    queryParams = { ...queryParams, search: uiStore.searchText };

    const query = this.appendQueryParams('people/posts', queryLimit, {
      ...queryParams,
      sortBy: 'created'
    });
    try {
      let ps = await this.fetchPeoplePosts(query);
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
          (n: any) => uiStore.setPeoplePostsPageNumber(n),
          queryParams
        );
      }
      return ps;
    } catch (e) {
      console.log('fetch failed getPeoplePosts: ', e);
      return [];
    }
  }

  @memo({
    resolver: (...args: any[]) => JSON.stringify({ args }),
    cache: new Map()
  })
  private async fetchPeoplePosts(query: string) {
    return await api.get(query);
  }

  @persist('list')
  peopleWanteds: PersonWanted[] = [];

  @action setPeopleWanteds(wanteds: PersonWanted[]) {
    this.peopleWanteds = wanteds;
  }

  getWantedsPrevParams?: QueryParams = {};

  async getPeopleWanteds(params?: QueryParams): Promise<PersonWanted[]> {
    const queryParams: QueryParams = {
      limit: 100,
      sortBy: 'created',
      search: uiStore.searchText ?? '',
      page: 1,
      resetPage: false,
      ...params
    };

    if (params) {
      // save previous params
      this.getWantedsPrevParams = queryParams;
    }

    // if we don't pass the params, we should use previous params for invalidate query
    const query2 = this.appendQueryParams(
      'bounty/all',
      queryLimit,
      params ? queryParams : this.getWantedsPrevParams
    );

    try {
      const ps2 = await api.get(query2);
      const ps3: any[] = [];

      if (ps2) {
        for (let i = 0; i < ps2.length; i++) {
          const bounty = { ...ps2[i].bounty };
          let assignee;
          let organization;
          const owner = { ...ps2[i].owner };

          if (bounty.assignee) {
            assignee = { ...ps2[i].assignee };
          }

          if (bounty.org_uuid) {
            organization = { ...ps2[i].organization };
          }

          ps3.push({
            body: { ...bounty, assignee: assignee || '' },
            person: { ...owner, wanteds: [] } || { wanteds: [] },
            organization: { ...organization }
          });
        }
      }

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.setPeopleWanteds(ps3);
        uiStore.setPeopleWantedsPageNumber(1);
      } else {
        // all other cases, merge
        const wanteds = this.doPageListMerger(
          this.peopleWanteds,
          ps3,
          (n: any) => uiStore.setPeopleWantedsPageNumber(n),
          queryParams,
          'wanted'
        );
        this.setPeopleWanteds(wanteds);
      }
      return ps3;
    } catch (e) {
      console.log('fetch failed getPeopleWanteds: ', e);
      return [];
    }
  }

  personAssignedWanteds: PersonWanted[] = [];

  @action setPersonWanteds(wanteds: PersonWanted[]) {
    this.personAssignedWanteds = wanteds;
  }

  async getPersonAssignedWanteds(queryParams?: any, pubkey?: string): Promise<PersonWanted[]> {
    queryParams = { ...queryParams, search: uiStore.searchText };

    const query = this.appendQueryParams(`people/wanteds/assigned/${pubkey}`, queryLimit, {
      ...queryParams,
      sortBy: 'created'
    });
    try {
      const ps2 = await api.get(query);
      const ps3: any[] = [];

      if (ps2 && ps2.length) {
        for (let i = 0; i < ps2.length; i++) {
          const bounty = { ...ps2[i].bounty };
          let assignee;
          let organization;
          const owner = { ...ps2[i].owner };

          if (bounty.assignee) {
            assignee = { ...ps2[i].assignee };
          }

          if (bounty.org_uuid) {
            organization = { ...ps2[i].organization };
          }

          ps3.push({
            body: { ...bounty, assignee: assignee || '' },
            person: { ...owner, wanteds: [] } || { wanteds: [] },
            organization: { ...organization }
          });
        }
      }

      return ps3;
    } catch (e) {
      console.log('fetch failed getPersonAssignedWanteds: ', e);
      return [];
    }
  }

  createdWanteds: PersonWanted[] = [];

  @action setCreatedWanteds(wanteds: PersonWanted[]) {
    this.createdWanteds = wanteds;
  }

  async getPersonCreatedWanteds(queryParams?: any, pubkey?: string): Promise<PersonWanted[]> {
    queryParams = { ...queryParams, search: uiStore.searchText };

    const query = this.appendQueryParams(`people/wanteds/created/${pubkey}`, queryLimit, {
      ...queryParams,
      sortBy: 'created'
    });
    try {
      const ps2 = await api.get(query);
      const ps3: any[] = [];

      if (ps2 && ps2.length) {
        for (let i = 0; i < ps2.length; i++) {
          const bounty = { ...ps2[i].bounty };
          let assignee;
          let organization;
          const owner = { ...ps2[i].owner };

          if (bounty.assignee) {
            assignee = { ...ps2[i].assignee };
          }

          if (bounty.org_uuid) {
            organization = { ...ps2[i].organization };
          }

          ps3.push({
            body: { ...bounty, assignee: assignee || '' },
            person: { ...owner, wanteds: [] } || { wanteds: [] },
            organization: { ...organization }
          });
        }
      }

      this.setCreatedWanteds(ps3);

      return ps3;
    } catch (e) {
      console.log('fetch failed getPersonCreatedWanteds: ', e);
      return [];
    }
  }

  async getWantedById(id: number): Promise<PersonWanted[]> {
    try {
      const ps2 = await api.get(`bounty/id/${id}`);
      const ps3: any[] = [];

      if (ps2 && ps2.length) {
        for (let i = 0; i < ps2.length; i++) {
          const bounty = { ...ps2[i].bounty };
          let assignee;
          let organization;
          const owner = { ...ps2[i].owner };

          if (bounty.assignee) {
            assignee = { ...ps2[i].assignee };
          }

          if (bounty.org_uuid) {
            organization = { ...ps2[i].organization };
          }

          ps3.push({
            body: { ...bounty, assignee: assignee || '' },
            person: { ...owner, wanteds: [] } || { wanteds: [] },
            organization: { ...organization }
          });
        }
      }

      return ps3;
    } catch (e) {
      console.log('fetch failed getWantedById: ', e);
      return [];
    }
  }

  async getOrganizationWanted(uuid: string, queryParams?: any): Promise<PersonWanted[]> {
    queryParams = { ...queryParams, search: uiStore.searchText };
    try {
      const ps2 = await api.get(`organizations/bounties/${uuid}`);
      const ps3: any[] = [];

      if (ps2 && ps2.length) {
        for (let i = 0; i < ps2.length; i++) {
          const bounty = { ...ps2[i].bounty };
          let assignee;
          let organization;
          const owner = { ...ps2[i].owner };

          if (bounty.assignee) {
            assignee = { ...ps2[i].assignee };
          }

          if (bounty.org_uuid) {
            organization = { ...ps2[i].organization };
          }

          ps3.push({
            body: { ...bounty, assignee: assignee || '' },
            person: { ...owner, wanteds: [] } || { wanteds: [] },
            organization: { ...organization }
          });
        }
      }

      // for search always reset page
      if (queryParams && queryParams.resetPage) {
        this.setPeopleWanteds(ps3);
        uiStore.setPeopleWantedsPageNumber(1);
      } else {
        // all other cases, merge
        const wanteds = this.doPageListMerger(
          this.peopleWanteds,
          ps3,
          (n: any) => uiStore.setPeopleWantedsPageNumber(n),
          queryParams,
          'wanted'
        );

        this.setPeopleWanteds(wanteds);
      }
      return ps3;
    } catch (e) {
      console.log('fetch failed getOrganizationWanted: ', e);
      return [];
    }
  }

  async getBountyCount(personKey: string, tabType: string): Promise<number> {
    try {
      const count = await api.get(`bounty/count/${personKey}/${tabType}`);

      return count;
    } catch (e) {
      console.log('fetch failed getCreatedWanteds: ', e);
      return 0;
    }
  }

  @persist('list')
  peopleOffers: PersonOffer[] = [];

  async getPeopleOffers(queryParams?: any): Promise<PersonOffer[]> {
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
          (n: any) => uiStore.setPeopleOffersPageNumber(n),
          queryParams
        );
      }

      return ps;
    } catch (e) {
      console.log('fetch failed getPeopleOffers: ', e);
      return [];
    }
  }

  doPageListMerger(
    currentList: any[],
    newList: any[],
    setPage: (any) => void,
    queryParams?: any,
    type?: string
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

    const set = new Set();
    if (type === 'wanted') {
      const uniqueArray = l.filter((item: any) => {
        if (item.body && item.body.id && !set.has(item.body.id)) {
          set.add(item.body.id);
          return true;
        }
        return false;
      }, set);
      return uniqueArray;
    }

    return l;
  }

  @persist('list')
  _activePerson: Person[] = [];

  get activePerson() {
    return this._activePerson.map((person: Person) => ({
      ...person,
      img: person.img || this.getUserAvatarPlaceholder(person.owner_pubkey)
    }));
  }

  set activePerson(p: Person[]) {
    this._activePerson = p;
  }

  setActivePerson(p: Person) {
    this.activePerson = [p];
  }

  @memo()
  async getPersonByPubkey(pubkey: string): Promise<Person> {
    const p = await api.get(`person/${pubkey}`);
    return p;
  }

  @memo()
  async getPersonById(id: number): Promise<Person> {
    const p = await api.get(`person/id/${id}`);
    this.setActivePerson(p);
    return p;
  }

  async getPersonByGithubName(github: string): Promise<Person> {
    const p = await api.get(`person/githubname/${github}`);
    return p;
  }

  // this method merges the relay self data with the db self data, they each hold different data

  async getSelf(me: any) {
    const self = me || uiStore.meInfo;
    if (self) {
      const p = await api.get(`person/${self.owner_pubkey}`);

      // get request for super_admin_array.
      const getSuperAdmin = async () => {
        try {
          const response = await api.get(`admin_pubkeys`);
          const admin_keys = response?.pubkeys;
          if (admin_keys !== null) {
            return !!admin_keys.find((value: any) => value === self.owner_pubkey);
          } else {
            return false;
          }
        } catch (error) {
          return false;
        }
      };

      const isSuperAdmin = await getSuperAdmin();
      const updateSelf = { ...self, ...p, isSuperAdmin: isSuperAdmin };
      uiStore.setMeInfo(updateSelf);
    }
  }

  async claimBadgeOnLiquid(body: ClaimOnLiquid): Promise<any> {
    try {
      const [r, error] = await this.doCallToRelay('POST', 'claim_on_liquid', body);
      if (error) throw error;
      if (!r) return; // tor user will return here

      return r;
    } catch (e) {
      console.log('Error claimBadgeOnLiquid: ', e);
    }
  }

  async sendBadgeOnLiquid(body: ClaimOnLiquid): Promise<any> {
    try {
      const [r, error] = await this.doCallToRelay('POST', 'claim_on_liquid', body);
      if (error) throw error;
      if (!r) return; // tor user will return here

      return r;
    } catch (e) {
      console.log('Error sendBadgeOnLiquid: ', e);
    }
  }

  async refreshJwt() {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;

      const r: any = await fetch(`${TribesURL}/refresh_jwt`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json',
          Accept: 'application/json'
        }
      });

      const j = await r.json();

      if (this.lnToken) {
        this.lnToken = j.jwt;
        return j;
      }

      return j;
    } catch (e) {
      console.log('Error refreshJwt: ', e);
      // could not refresh jwt, logout!
      return null;
    }
  }

  async getUsdToSatsExchangeRate() {
    try {
      // get rate for 1 USD
      const res: any = await fetch('https://blockchain.info/tobtc?currency=USD&value=1', {
        method: 'GET'
      });
      const j = await res.json();
      // 1 bitcoin is 1 million satoshis
      const satoshisInABitcoin = 0.00000001;
      const exchangeRate = j / satoshisInABitcoin;

      uiStore.setUsdToSatsExchangeRate(exchangeRate);

      return exchangeRate;
    } catch (e) {
      console.log('Error getUsdToSatsExchangeRate: ', e);
      // could not refresh jwt, logout!
      return null;
    }
  }

  async deleteProfile() {
    try {
      const info = uiStore.meInfo;
      let request = 'profile';
      if (this.lnToken) request = `person/${info?.id}`;

      const [r, error] = await this.doCallToRelay('DELETE', request, info);
      if (error) throw error;
      if (!r) return; // tor user will return here

      uiStore.setMeInfo(null);
      uiStore.setSelectingPerson(0);
      uiStore.setSelectedPerson(0);

      const j = await r.json();
      return j;
    } catch (e) {
      console.log('Error deleteProfile: ', e);
      // could not delete profile!
      return null;
    }
  }

  async saveProfile(body: any) {
    if (!uiStore.meInfo) return null;
    const info = uiStore.meInfo;
    if (!body) return; // avoid saving bad state

    if (body.price_to_meet) body.price_to_meet = parseInt(body.price_to_meet); // must be an int

    try {
      if (this.lnToken) {
        const r = await fetch(`${TribesURL}/person`, {
          method: 'POST',
          body: JSON.stringify({
            ...body
          }),
          mode: 'cors',
          headers: {
            'x-jwt': info.tribe_jwt,
            'Content-Type': 'application/json'
          }
        });
        if (!r) return;
        // first time profile makers will need this on first login
        if (r.status === 200) {
          const p = await r.json();
          const updateSelf = { ...info, ...p };
          uiStore.setMeInfo(updateSelf);
        }
      } else {
        const [r, error] = await this.doCallToRelay('POST', 'profile', body);
        if (error) throw error;
        if (!r) return;

        // first time profile makers will need this on first login
        if (!body.id) {
          const j = await r.json();
          if (j.response && j.response.id) {
            body.id = j.response.id;
          }
        }

        const updateSelf = { ...info, ...body };
        await this.getSelf(updateSelf);

        uiStore.setToasts([
          {
            id: '1',
            title: 'Saved.'
          }
        ]);
      }
    } catch (e) {
      console.log('Error saveProfile: ', e);
    }
  }

  async saveBounty(body: any): Promise<void> {
    const info = uiStore.meInfo as any;
    if (!info && !body) {
      console.log('Youre not logged in');
      return;
    }

    if (!body.coding_languages || !body.coding_languages.length) {
      body.coding_languages = [];
    } else {
      const languages: string[] = [];
      body.coding_languages.forEach((lang: any) => {
        languages.push(lang.value);
      });

      body.coding_languages = languages;
    }

    try {
      const request = `bounty?token=${info?.tribe_jwt}`;
      //TODO: add some sort of authentication
      const response = await fetch(`${TribesURL}/${request}`, {
        method: 'POST',
        body: JSON.stringify({
          ...body
        }),
        mode: 'cors',
        headers: {
          'x-jwt': info?.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      if (response.status) {
        this.getPeopleWanteds({ resetPage: true });
      }
      return;
    } catch (e) {
      console.log(e);
    }
  }

  async deleteBounty(created: number, owner_pubkey: string): Promise<void> {
    const info = uiStore.meInfo as any;
    if (!info) {
      console.log('Youre not logged in');
      return;
    }

    try {
      const request = `bounty/${owner_pubkey}/${created}`;
      //TODO: add some sort of authentication
      const response = await fetch(`${TribesURL}/${request}`, {
        method: 'DELETE',
        mode: 'cors',
        headers: {
          'x-jwt': info?.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });
      if (response.status) {
        await this.getPeopleWanteds({ resetPage: true });
      }
      return;
    } catch (e) {
      console.log(e);
    }
  }

  // this method is used whenever changing data from the frontend,
  // forks between tor users and non-tor
  async doCallToRelay(method: string, path: string, body: any): Promise<any> {
    let error: any = null;

    const info = uiStore.meInfo as any;
    const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`;
    if (!info) {
      error = new Error('Youre not logged in');
      return [null, error];
    }

    if (this.lnToken) {
      const response = await fetch(`${URL}/${path}`, {
        method: method,
        body: JSON.stringify({
          ...body
        }),
        mode: 'cors',
        headers: {
          'x-jwt': info.jwt,
          'Content-Type': 'application/json'
        }
      });

      return [response, error];
    } else {
      // fork between tor users non authentiacted and not
      if (this.isTorSave() || info.url.startsWith('http://')) {
        this.submitFormViaApp(method, path, body);
        return [null, null];
      }

      const response = await fetch(`${URL}/${path}`, {
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
  }

  async submitFormViaApp(method: string, path: string, body: any) {
    try {
      const torSaveURL = await this.getTorSaveURL(method, path, body);
      uiStore.setTorFormBodyQR(torSaveURL);
    } catch (e) {
      console.log('Error submitFormViaApp: ', e);
    }
  }

  async setExtrasPropertyAndSave(
    extrasName: string,
    propertyName: string,
    created: number,
    newPropertyValue: any
  ): Promise<any> {
    if (uiStore.meInfo) {
      const clonedMeInfo = { ...uiStore.meInfo };
      const clonedExtras = clonedMeInfo?.extras;
      const clonedEx: any = clonedExtras && clonedExtras[extrasName];
      const targetIndex = clonedEx?.findIndex((f: any) => f.created === created);

      if (clonedEx && (targetIndex || targetIndex === 0) && targetIndex > -1) {
        try {
          clonedEx[targetIndex][propertyName] = newPropertyValue;
          clonedMeInfo.extras[extrasName] = clonedEx;
          await this.saveProfile(clonedMeInfo);
          return [clonedEx, targetIndex];
        } catch (e) {
          console.log('Error setExtrasPropertyAndSave', e);
        }
      }

      return [null, null];
    }
  }

  // function to update many value in wanted array of object
  async setExtrasMultipleProperty(
    dataObject: object,
    extrasName: string,
    created: number
  ): Promise<any> {
    if (uiStore.meInfo) {
      const clonedMeInfo = { ...uiStore.meInfo };
      const clonedExtras = clonedMeInfo?.extras;
      const clonedEx: any = clonedExtras && clonedExtras[extrasName];
      const targetIndex = clonedEx?.findIndex((f: any) => f.created === created);

      if (clonedEx && (targetIndex || targetIndex === 0) && targetIndex > -1) {
        try {
          clonedEx[targetIndex] = { ...clonedEx?.[targetIndex], ...dataObject };
          clonedMeInfo.extras[extrasName] = clonedEx;
          await this.saveProfile(clonedMeInfo);
          return [clonedEx, targetIndex];
        } catch (e) {
          console.log('Error setExtrasMultipleProperty', e);
        }
      }

      return [null, null];
    }
  }

  async deleteFavorite() {
    const body: any = {};

    if (!body) return; // avoid saving bad state

    const info = uiStore.meInfo as any;
    if (!info) return;
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
          'x-jwt': info.tribe_jwt,
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
      console.log('Error deleteFavorite', e);
    }
  }

  async getBountyHeaderData() {
    try {
      const data = await api.get('people/wanteds/header');
      return data;
    } catch (e) {
      console.log('Error getBountyHeaderData', e);
      return '';
    }
  }

  @observable
  lnauth: LnAuthData = { encode: '', k1: '' };

  @action setLnAuth(lnData: LnAuthData) {
    this.lnauth = lnData;
  }

  @persist('object')
  @observable
  lnToken = '';

  @action setLnToken(token: string) {
    this.lnToken = token;
  }

  @action async getLnAuth(): Promise<LnAuthData> {
    try {
      const data = await api.get(`lnauth?socketKey=${uiStore.websocketToken}`);
      this.setLnAuth(data);
      return data;
    } catch (e) {
      console.log('Error getLnAuth', e);
      return { encode: '', k1: '' };
    }
  }

  @observable
  lnInvoice = '';

  @action setLnInvoice(invoice: string) {
    this.lnInvoice = invoice;
  }

  async getLnInvoice(body: {
    amount: number;
    memo: string;
    owner_pubkey: string;
    user_pubkey: string;
    created: string;
    type: 'KEYSEND' | 'ASSIGN';
    assigned_hours?: number;
    commitment_fee?: number;
    bounty_expires?: string;
    route_hint?: string;
  }): Promise<LnInvoice> {
    try {
      const data = await api.post(
        'invoices',
        {
          amount: body.amount.toString(),
          memo: body.memo,
          owner_pubkey: body.owner_pubkey,
          user_pubkey: body.user_pubkey,
          created: body.created,
          type: body.type,
          assigned_hours: body.assigned_hours,
          commitment_fee: body.commitment_fee,
          bounty_expires: body.bounty_expires,
          websocket_token: uiStore.meInfo?.websocketToken,
          route_hint: body.route_hint
        },
        {
          'Content-Type': 'application/json'
        }
      );
      if (data.success) {
        this.setLnInvoice(data.response.invoice);
      }
      return data;
    } catch (e) {
      console.log('Error getLnInvoice', e);
      return { success: false, response: { invoice: '' } };
    }
  }

  async getBudgetInvoice(body: {
    amount: number;
    org_uuid: string;
    sender_pubkey: string;
    websocket_token: string;
    payment_type: string;
  }): Promise<LnInvoice> {
    try {
      const data = await api.post(
        'budgetinvoices',
        {
          amount: body.amount,
          org_uuid: body.org_uuid,
          sender_pubkey: body.sender_pubkey,
          websocket_token: body.websocket_token,
          payment_type: body.payment_type
        },
        {
          'Content-Type': 'application/json'
        }
      );
      if (data.success) {
        this.setLnInvoice(data.response.invoice);
      }
      return data;
    } catch (e) {
      return { success: false, response: { invoice: '' } };
    }
  }

  @action async deleteBountyAssignee(body: {
    owner_pubkey: string;
    created: string;
  }): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/bounty/assignee`, {
        method: 'DELETE',
        mode: 'cors',
        body: JSON.stringify({
          ...body
        }),
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r;
    } catch (e) {
      console.log('Error deleteBountyAssignee', e);
      return false;
    }
  }

  @observable
  organizations: Organization[] = [];

  @action setOrganizations(organizations: Organization[]) {
    this.organizations = organizations;
  }

  @action async getUserOrganizations(id: number): Promise<Organization[]> {
    try {
      if (!uiStore.meInfo) return [];
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/user/${id}`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      const data = await r.json();
      this.setOrganizations(data);
      return await data;
    } catch (e) {
      console.log('Error getUserOrganizations', e);
      return [];
    }
  }

  @action async addOrganization(body: { name: string; img: string }): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations`, {
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify({
          ...body
        }),
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r;
    } catch (e) {
      console.log('Error addOrganization', e);
      return false;
    }
  }

  async getOrganizationUsersCount(uuid: string): Promise<number> {
    try {
      const r: any = await fetch(`${TribesURL}/organizations/users/${uuid}/count`, {
        method: 'GET',
        mode: 'cors'
      });

      return r.json();
    } catch (e) {
      console.log('Error getOrganizationUsersCount', e);
      return 0;
    }
  }

  async getOrganizationUsers(uuid: string): Promise<Person[]> {
    try {
      const r: any = await fetch(`${TribesURL}/organizations/users/${uuid}`, {
        method: 'GET',
        mode: 'cors'
      });

      return r.json();
    } catch (e) {
      console.log('Error getOrganizationUsers', e);
      return [];
    }
  }

  @action async addOrganizationUser(body: {
    owner_pubkey: string;
    org_uuid: string;
  }): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/users/${body.org_uuid}`, {
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify({
          ...body
        }),
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r;
    } catch (e) {
      console.log('Error addOrganizationUser', e);
      return false;
    }
  }

  @action async deleteOrganizationUser(body: any, uuid: string): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/users/${uuid}`, {
        method: 'DELETE',
        mode: 'cors',
        body: JSON.stringify({
          ...body
        }),
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r;
    } catch (e) {
      console.log('Error deleteOrganizationUser', e);
      return false;
    }
  }

  @observable
  bountyRoles: BountyRoles[] = [];

  @action setBountyRoles(roles: BountyRoles[]) {
    this.bountyRoles = roles;
  }

  async getRoles(): Promise<BountyRoles[]> {
    try {
      if (!uiStore.meInfo) return [];
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/bounty/roles`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      const roles = await r.json();
      this.setBountyRoles(roles);

      return roles;
    } catch (e) {
      console.log('Error getRoles', e);
      return [];
    }
  }

  async getUserRoles(uuid: string, user: string): Promise<any[]> {
    try {
      if (!uiStore.meInfo) return [];
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/users/role/${uuid}/${user}`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r.json();
    } catch (e) {
      console.log('Error getUserRoles', e);
      return [];
    }
  }

  async addUserRoles(body: any, uuid: string, user: string): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/users/role/${uuid}/${user}`, {
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify(body),
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r;
    } catch (e) {
      console.log('Error addUserRoles', e);
      return false;
    }
  }

  async updateBountyPaymentStatus(created: number): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/bounty/paymentstatus/${created}`, {
        method: 'POST',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r;
    } catch (e) {
      console.log('Error updateBountyPaymentStatus', e);
      return false;
    }
  }

  async getOrganizationBudget(uuid: string): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/budget/${uuid}`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r.json();
    } catch (e) {
      console.log('Error getOrganizationBudget', e);
      return false;
    }
  }

  async makeBountyPayment(body: {
    id: number;
    receiver_pubkey: string;
    websocket_token: string;
    route_hint: string;
  }): Promise<any> {
    try {
      if (!uiStore.meInfo) return null;
      const info = uiStore.meInfo;

      const r: any = await fetch(`${TribesURL}/bounty/pay/${body.id}`, {
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify(body),
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r;
    } catch (e) {
      console.log('Error makeBountyPayment', e);
      return false;
    }
  }

  async getPaymentHistories(uuid: string): Promise<PaymentHistory[]> {
    try {
      if (!uiStore.meInfo) return [];
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/payments/${uuid}`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r.json();
    } catch (e) {
      console.log('Error getPaymentHistories', e);
      return [];
    }
  }

  async getBudgettHistories(uuid: string): Promise<BudgetHistory[]> {
    try {
      if (!uiStore.meInfo) return [];
      const info = uiStore.meInfo;
      const r: any = await fetch(`${TribesURL}/organizations/budget/history/${uuid}`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'x-jwt': info.tribe_jwt,
          'Content-Type': 'application/json'
        }
      });

      return r.json();
    } catch (e) {
      console.log('Error gettHistories', e);
      return [];
    }
  }
}

export const mainStore = new MainStore();
