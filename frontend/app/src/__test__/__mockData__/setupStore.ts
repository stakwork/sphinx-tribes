import { mainStore } from '../../store/main';
import { uiStore } from '../../store/ui';
import { person } from './persons';
import { user } from './user';
import { userAssignedTickets } from './userTickets';

export const setupStore = () => {
  mainStore.setPeople([person]);
  uiStore.setMeInfo(user);
  mainStore.setPersonWanteds(userAssignedTickets);
  uiStore.setSelectedPerson(user.id as number);
};
