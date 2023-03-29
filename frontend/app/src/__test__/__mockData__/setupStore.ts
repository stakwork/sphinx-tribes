import { mainStore } from '../../store/main';
import { uiStore } from '../../store/ui';
import { person } from './persons';
import { user } from './user';

export const setupStore = () => {
  mainStore.setPeople([person]);
  uiStore.setMeInfo(user);
  uiStore.setSelectedPerson(user.id as number);
};
