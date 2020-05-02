import { observable, action } from 'mobx'

class UiStore {
  @observable ready: boolean = false
  @action setReady(ready:boolean){
    this.ready = ready
  }
}

export const uiStore = new UiStore()