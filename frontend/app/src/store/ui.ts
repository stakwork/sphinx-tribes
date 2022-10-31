import { observable, action } from 'mobx';
import { persist } from 'mobx-persist';
import tags from '../tribes/tags';
import { Extras } from '../form/inputs/widgets/interfaces';

const tagLabels = Object.keys(tags);
const initialTags = tagLabels.map((label) => {
  return <EuiSelectableOption>{ label };
});

export type EuiSelectableOptionCheckedType = 'on' | 'off' | undefined;

export interface EuiSelectableOption {
  label: string;
  checked?: EuiSelectableOptionCheckedType;
}

class UiStore {
  @observable ready: boolean = false;
  @action setReady(ready: boolean) {
    this.ready = ready;
  }

  @observable tags: EuiSelectableOption[] = initialTags;
  @action setTags(t: EuiSelectableOption[]) {
    this.tags = t;
  }

  @observable searchText: string = '';
  @action setSearchText(s: string) {
    this.searchText = s.toLowerCase();
  }

  @observable usdToSatsExchangeRate: number = 0;
  @action setUsdToSatsExchangeRate(n: number) {
    this.usdToSatsExchangeRate = n;
  }

  @observable editMe: boolean = false;
  @action setEditMe(b: boolean) {
    this.editMe = b;
  }

  @observable peoplePageNumber: number = 1;
  @action setPeoplePageNumber(n: number) {
    this.peoplePageNumber = n;
  }

  @observable peoplePostsPageNumber: number = 1;
  @action setPeoplePostsPageNumber(n: number) {
    this.peoplePostsPageNumber = n;
  }

  @observable peopleWantedsPageNumber: number = 1;
  @action setPeopleWantedsPageNumber(n: number) {
    this.peopleWantedsPageNumber = n;
  }

  @observable peopleOffersPageNumber: number = 1;
  @action setPeopleOffersPageNumber(n: number) {
    this.peopleOffersPageNumber = n;
  }

  @observable tribesPageNumber: number = 1;
  @action setTribesPageNumber(n: number) {
    this.tribesPageNumber = n;
  }

  @observable selectedPerson: number = 0;
  @action setSelectedPerson(n: number) {
    this.selectedPerson = n;
  }

  // this is for animations, if you deselect as a component is fading out,
  // it empties and looks broke for a second
  @observable selectingPerson: number = 0;
  @action setSelectingPerson(n: number) {
    this.selectingPerson = n;
  }

  @observable selectedBot: string = '';
  @action setSelectedBot(n: string) {
    this.selectedBot = n;
  }

  // this is for animations, if you deselect as a component is fading out,
  // it empties and looks broke for a second
  @observable selectingBot: string = '';
  @action setSelectingBot(n: string) {
    this.selectingBot = n;
  }

  @observable toasts: any = [];
  @action setToasts(n: any) {
    this.toasts = n;
  }

  @observable personViewOpenTab: string = '';
  @action setPersonViewOpenTab(s: string) {
    this.personViewOpenTab = s;
  }

  @observable lastGithubRepo: string = '';
  @action setLastGithubRepo(s: string) {
    this.lastGithubRepo = s;
  }

  @observable torFormBodyQR: string = '';
  @action setTorFormBodyQR(s: string) {
    this.torFormBodyQR = s;
  }

  @observable openGithubIssues: any = [];
  @action setOpenGithubIssues(a: any) {
    this.openGithubIssues = a;
  }

  @observable badgeList: any = [];
  @action setBadgeList(a: any) {
    this.badgeList = a;
  }

  @observable language: string = '';
  @action setLanguage(s: string) {
    this.language = s;
  }

  @persist('object') @observable meInfo: MeData = null;
  @action setMeInfo(t: MeData) {
    if (t) {
      if (t.photo_url && !t.img) t.img = t.photo_url;
      if (!t.owner_alias) t.owner_alias = t.alias;
      if (!t.owner_pubkey) t.owner_pubkey = t.pubkey;
    }
    this.meInfo = t;
  }

  @observable showSignIn: boolean = false;
  @action setShowSignIn(b: boolean) {
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
