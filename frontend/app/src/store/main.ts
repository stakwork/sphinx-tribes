import { observable, action } from 'mobx'

export class MainStore {
  @observable
  hi: string = "hello"
}

export const mainStore = new MainStore()