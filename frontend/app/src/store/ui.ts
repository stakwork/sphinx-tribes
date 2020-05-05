import { observable, action } from 'mobx'
import tags from '../components/tags'

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
}

export const uiStore = new UiStore()