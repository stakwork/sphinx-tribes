import { mainStore } from '../../store/main';
import { uiStore } from '../../store/ui';
import { person, userTickets } from './persons';
import { user } from './user';

export const setupStore = () => {
  mainStore.setPeople([person]);
  uiStore.setMeInfo(user);
  mainStore.setPersonWanteds(userTickets);
  uiStore.setSelectedPerson(user.id as number);
};
