import { mainStore } from '../../store/main';
import { uiStore } from '../../store/ui';
import { person } from './persons';
import { user } from './user';
import { userAssignedBounties } from './userTickets';

export const setupStore = () => {
  mainStore.setPeople([person]);
  uiStore.setMeInfo(user);
  mainStore.setPersonBounties(userAssignedBounties);
  uiStore.setSelectedPerson(user.id as number);
};
