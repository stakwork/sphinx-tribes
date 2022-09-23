import React from 'react';
import { uiStore } from './ui';
import { mainStore } from './main';
import { create } from 'mobx-persist';

const hydrate = create({ storage: localStorage });

Promise.all([hydrate('main', mainStore), hydrate('ui', uiStore)]).then(() => {
  console.log('useStore ready!');
  uiStore.setReady(true);
});

const ctx = React.createContext({
  ui: uiStore,
  main: mainStore
});

export const useStores = () => React.useContext(ctx);
