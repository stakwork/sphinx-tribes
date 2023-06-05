import { makeAutoObservable } from 'mobx';
import { persist } from 'mobx-persist';
import { Extras } from '../components/form/inputs/widgets/interfaces';
import tags from '../tribes/tags';

const tagLabels = Object.keys(tags);
const initialTags = tagLabels.map((label) => ({ label } as EuiSelectableOption));

export type EuiSelectableOptionCheckedType = 'on' | 'off' | undefined;

export interface EuiSelectableOption {
  label: string;
  checked?: EuiSelectableOptionCheckedType;
}

class UiStore {
  ready = false;

  constructor() {
    makeAutoObservable(this);
  }

  setReady(ready: boolean) {
    this.ready = ready;
  }

  tags: EuiSelectableOption[] = initialTags;
  setTags(t: EuiSelectableOption[]) {
    this.tags = t;
  }

  searchText = '';
  setSearchText(s: string) {
    this.searchText = s.toLowerCase();
  }

  usdToSatsExchangeRate = 0;
  setUsdToSatsExchangeRate(n: number) {
    this.usdToSatsExchangeRate = n;
  }

  editMe = false;
  setEditMe(b: boolean) {
    this.editMe = b;
  }

  peoplePageNumber = 1;
  setPeoplePageNumber(n: number) {
    this.peoplePageNumber = n;
  }

  peoplePostsPageNumber = 1;
  setPeoplePostsPageNumber(n: number) {
    this.peoplePostsPageNumber = n;
  }

  peopleWantedsPageNumber = 1;
  setPeopleWantedsPageNumber(n: number) {
    this.peopleWantedsPageNumber = n;
  }

  peopleOffersPageNumber = 1;
  setPeopleOffersPageNumber(n: number) {
    this.peopleOffersPageNumber = n;
  }

  tribesPageNumber = 1;
  setTribesPageNumber(n: number) {
    this.tribesPageNumber = n;
  }

  selectedPerson = 0;
  setSelectedPerson(n: number | undefined) {
    if (n) this.selectedPerson = n;
  }

  // this is for animations, if you deselect as a component is fading out,
  // it empties and looks broke for a second
  selectingPerson = 0;
  setSelectingPerson(n: number | undefined) {
    if (n) this.selectingPerson = n;
  }

  selectedBot = '';
  setSelectedBot(n: string) {
    this.selectedBot = n;
  }

  // this is for animations, if you deselect as a component is fading out,
  // it empties and looks broke for a second
  selectingBot = '';
  setSelectingBot(n: string) {
    this.selectingBot = n;
  }

  toasts: any = [];
  setToasts(n: any) {
    this.toasts = n;
  }

  personViewOpenTab = '';
  setPersonViewOpenTab(s: string) {
    this.personViewOpenTab = s;
  }

  lastGithubRepo = '';
  setLastGithubRepo(s: string) {
    this.lastGithubRepo = s;
  }

  torFormBodyQR = '';
  setTorFormBodyQR(s: string) {
    this.torFormBodyQR = s;
  }

  openGithubIssues: any = [];
  setOpenGithubIssues(a: any) {
    this.openGithubIssues = a;
  }

  badgeList: any = [];
  setBadgeList(a: any) {
    this.badgeList = a;
  }

  language = '';
  setLanguage(s: string) {
    this.language = s;
  }

  @persist('object') meInfo: MeData = null;
  setMeInfo(t: MeData) {
    if (t) {
      if (t.photo_url && !t.img) t.img = t.photo_url;
      if (!t.owner_alias) t.owner_alias = t.alias;
      if (!t.owner_pubkey) t.owner_pubkey = t.pubkey;
    }
    this.meInfo = t;
  }

  @persist('object') connection_string = '';
  setConnectionString(code: string) {
    this.connection_string = code;
  }

  showSignIn = false;
  setShowSignIn(b: boolean) {
    this.showSignIn = b;
  }
}

export type MeData = MeInfo | null;

export interface MeInfo {
  id?: number;
  pubkey: string;
  owner_pubkey?: string;
  photo_url: string;
  alias: string;
  img?: string;
  owner_alias?: string;
  github_issues?: any[];
  route_hint: string;
  contact_key: string;
  price_to_meet: number;
  jwt: string;
  url: string;
  description: string;
  verification_signature: string;
  twitter_confirmed?: boolean;
  extras: Extras;
  isSuperAdmin: boolean;
  websocketToken?: string
}

export const emptyMeData: MeData = {
  pubkey: '',
  alias: '',
  route_hint: '',
  contact_key: '',
  price_to_meet: 0,
  photo_url: '',
  url: '',
  jwt: '',
  description: '',
  verification_signature: '',
  extras: {},
  isSuperAdmin: false
};

export const emptyMeInfo: MeInfo = {
  pubkey: '',
  alias: '',
  route_hint: '',
  contact_key: '',
  price_to_meet: 0,
  photo_url: '',
  url: '',
  jwt: '',
  description: '',
  verification_signature: '',
  extras: {},
  isSuperAdmin: false
};

export const uiStore = new UiStore();

// const hydrate = createPersist()
// hydrate('some', uiStore).then(() => { })
