import { observable, action } from 'mobx'
import tags from '../tribes/tags'

const tagLabels = Object.keys(tags)
const initialTags = tagLabels.map(label=>{
  return <EuiSelectableOption>{label}
})

export type EuiSelectableOptionCheckedType = 'on' | 'off' | undefined;

export interface EuiSelectableOption {
  label: string;
  checked?: EuiSelectableOptionCheckedType;
}

class UiStore {
  @observable ready: boolean = false
  @action setReady(ready:boolean){
    this.ready = ready
  }

  @observable tags: EuiSelectableOption[] = initialTags
  @action setTags(t:EuiSelectableOption[]){
    this.tags = t
  }

  @observable searchText: string = ''
  @action setSearchText(s:string){
    this.searchText = s
  }

  @observable editMe: boolean = true
  @action setEditMe(b:boolean){
    this.editMe = b
  }

  @observable tokens: TokensData = null
  @action setTokens(t:Tokens){
    this.tokens = t
  }
}

export type TokensData = Tokens | null
export interface Tokens {
  pubkey: string
  memeToken: string
  tribesTokens: string
}

export const uiStore = new UiStore()