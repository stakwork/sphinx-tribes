import { makeAutoObservable } from 'mobx';

class ModalsVisibilityStore {
  constructor() {
    makeAutoObservable(this);
  }

  userEditModal = false;
  setUserEditModal(v: boolean) {
    this.userEditModal = v;
  }
}

export const modalsVisibilityStore = new ModalsVisibilityStore();
