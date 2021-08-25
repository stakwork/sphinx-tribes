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

  // constructor() {
  //   this.meInfo = emptyMeData
  // }

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

  @persist('object') @observable meInfo: MeData = null
  @action setMeInfo(t: MeData) {
    if (t) {
      t.img = t.photo_url
      if (!t.owner_alias) t.owner_alias = t.alias
    }
    this.meInfo = t
  }

  @persist('object') @observable challenge: Challenge = { challenge: '' }
  @action setChallenge(t: string) {
    let c = {
      challenge: t
    }
    this.challenge = c
  }
}

export type MeData = MeInfo | null

export interface Challenge {
  challenge: string
}
export interface MeInfo {
  id?: number
  pubkey: string
  photo_url: string
  alias: string
  img?: string
  owner_alias?: string
  route_hint: string
  contact_key: string
  price_to_meet: number
  jwt: string
  url: string
  description: string
  verification_signature: string
  extras: Extras,
  challenge?: string
}
export const emptyMeData: MeData = { pubkey: '', alias: '', route_hint: '', contact_key: '', price_to_meet: 0, photo_url: '', url: '', jwt: '', description: '', verification_signature: '', extras: {} }
export const emptyMeInfo: MeInfo = { pubkey: '', alias: '', route_hint: '', contact_key: '', price_to_meet: 0, photo_url: '', url: '', jwt: '', description: '', verification_signature: '', extras: {} }

export const uiStore = new UiStore()

const hydrate = createPersist()
hydrate('some', uiStore).then(() => {})
