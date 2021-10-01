import { observable, action } from 'mobx'
import { create as createPersist, persist } from 'mobx-persist'
import tags from '../tribes/tags'
import { Extras } from '../form/inputs/widgets/interfaces'

const tagLabels = Object.keys(tags)
const initialTags = tagLabels.map(label => {
  return <EuiSelectableOption>{ label }
})

export type EuiSelectableOptionCheckedType = 'on' | 'off' | undefined;

export interface EuiSelectableOption {
  label: string;
  checked?: EuiSelectableOptionCheckedType;
}

class UiStore {
  @observable ready: boolean = false
  @action setReady(ready: boolean) {
    this.ready = ready
  }

  @observable tags: EuiSelectableOption[] = initialTags
  @action setTags(t: EuiSelectableOption[]) {
    this.tags = t
  }

  @observable searchText: string = ''
  @action setSearchText(s: string) {
    this.searchText = s
  }

  @observable editMe: boolean = false
  @action setEditMe(b: boolean) {
    this.editMe = b
  }

  @observable selectedPerson: number = 0
  @action setSelectedPerson(n: number) {
    this.selectedPerson = n
  }

  // this is for animations, if you deselect as a component is fading out, 
  // it empties and looks broke for a second
  @observable selectingPerson: number = 0
  @action setSelectingPerson(n: number) {
    this.selectingPerson = n
  }

  @observable selectedBot: string = ''
  @action setSelectedBot(n: string) {
    this.selectedBot = n
  }

  // this is for animations, if you deselect as a component is fading out, 
  // it empties and looks broke for a second
  @observable selectingBot: string = ''
  @action setSelectingBot(n: string) {
    this.selectingBot = n
  }

  @observable toasts: any = []
  @action setToasts(n: any) {
    console.log('set toasts', n)
    this.toasts = n
  }

  @observable personViewOpenTab: string = ''
  @action setPersonViewOpenTab(s: string) {
    console.log('set setPersonViewOpenTab', s)
    this.personViewOpenTab = s
  }

  @observable lastGithubRepo: string = ''
  @action setLastGithubRepo(s: string) {
    console.log('set setLastGithubRepo', s)
    this.lastGithubRepo = s
  }



  @persist('object') @observable meInfo: MeData = null
  @action setMeInfo(t: MeData) {
    if (t) {
      if (t.photo_url && !t.img) t.img = t.photo_url
      if (!t.owner_alias) t.owner_alias = t.alias
      if (!t.owner_pubkey) t.owner_pubkey = t.pubkey
    }
    this.meInfo = t
  }

  @observable showSignIn: boolean = false
  @action setShowSignIn(b: boolean) {
    this.showSignIn = b
  }

}

export type MeData = MeInfo | null

export interface MeInfo {
  id?: number
  pubkey: string
  owner_pubkey?: string
  photo_url: string
  alias: string
  img?: string
  owner_alias?: string
  github_issues?: any[]
  route_hint: string
  contact_key: string
  price_to_meet: number
  jwt: string
  url: string
  description: string
  verification_signature: string
  extras: Extras
}
export const emptyMeData: MeData = { pubkey: '', alias: '', route_hint: '', contact_key: '', price_to_meet: 0, photo_url: '', url: '', jwt: '', description: '', verification_signature: '', extras: {} }
export const emptyMeInfo: MeInfo = { pubkey: '', alias: '', route_hint: '', contact_key: '', price_to_meet: 0, photo_url: '', url: '', jwt: '', description: '', verification_signature: '', extras: {} }

export const uiStore = new UiStore()

// const hydrate = createPersist()
// hydrate('some', uiStore).then(() => { })
