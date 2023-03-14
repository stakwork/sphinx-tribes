import React from 'react';
import { uiStore } from './ui';
import { mainStore } from './main';
import { create } from 'mobx-persist';
import { appEnv } from '../config/env';
export { modalsVisibility } from './modals'

(() => {
  if (appEnv.isTests) {
    return;
  }
  const hydrate = create({ storage: localStorage });

  Promise.all([hydrate('main', mainStore), hydrate('ui', uiStore)]).then(() => {
    uiStore.setReady(true);
  });
})();

const Context = React.createContext({
  ui: uiStore,
  main: mainStore
});

export const useStores = () => React.useContext(Context);

