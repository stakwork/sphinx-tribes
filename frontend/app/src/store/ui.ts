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

  @observable meInfo: MeData = null
  @action setMeInfo(t:MeData){
    this.meInfo = t
  }
}

export type MeData = MeInfo | null
export interface MeInfo {
  pubkey: string
  memeToken: string
  tribesToken: string
  photoUrl: string
  alias: string
  routeHint: string
  contactKey: string
}
const emptyMeData:MeData = {pubkey:'asdf',alias:'evan',memeToken:'',tribesToken:'',routeHint:'',contactKey:'',photoUrl:''}

export const uiStore = new UiStore()