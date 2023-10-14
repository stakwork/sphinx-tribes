import { makeAutoObservable } from 'mobx';

class ModalsVisibilityStore {
  constructor() {
    makeAutoObservable(this);
  }

  userEditModal = false;
  setUserEditModal(v: boolean) {
    this.userEditModal = v;
  }

  startupModal = false;
  setStartupModal(v: boolean) {
    this.startupModal = v;
  }
}

export const modalsVisibilityStore = new ModalsVisibilityStore();
