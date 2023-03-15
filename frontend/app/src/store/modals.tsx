import { observable, action } from 'mobx';
import React from 'react';

class ModalsVisibilityStore {

  @observable
  userEditModal = false;

  @action setUserEditModal(v: boolean) {
    this.userEditModal = v;
  }
}

export const modalsVisibility = new ModalsVisibilityStore();

const Context = React.createContext({} as ModalsVisibilityStore);
export const useModalsVisibility = () => React.useContext(Context);


export const WithModalStore = ({children}) => {
  return (
  <Context.Provider value={modalsVisibility}> 
    {children}
  </Context.Provider>
  );
}