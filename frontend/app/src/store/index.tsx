import React from 'react';
import { create } from 'mobx-persist';
import { configure } from 'mobx';
import { leaderboardStore } from 'leaderboard';
import { appEnv } from '../config/env';
import { uiStore } from './ui';
import { mainStore } from './main';
import { modalsVisibilityStore } from './modals';

(() => {
  if (appEnv.isTests) {
    return;
  }
  const hydrate = create({ storage: localStorage });

  Promise.all([hydrate('main', mainStore), hydrate('ui', uiStore)]).then(() => {
    uiStore.setReady(true);
  });
})();

configure({});
const Context = React.createContext({
  ui: uiStore,
  main: mainStore,
  modals: modalsVisibilityStore,
  leaderboard: leaderboardStore
});

export const WithStores = ({ children }: any) => (
  <Context.Provider
    value={{
      ui: uiStore,
      main: mainStore,
      modals: modalsVisibilityStore,
      leaderboard: leaderboardStore
    }}
  >
    {children}
  </Context.Provider>
);
export const useStores = () => React.useContext(Context);
