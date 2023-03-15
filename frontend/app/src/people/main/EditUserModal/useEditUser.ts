import { usePerson } from 'hooks';
import { useStores } from 'store';
import { useModalsVisibility } from 'store/modals';

export const useUserEdit = () => {
  const {ui} = useStores()
  const modals = useModalsVisibility();
  const personId = ui.selectedPerson;
  // const personId = useParams<{personId: string}>();
  
  const {canEdit, person} = usePerson(Number(personId));

  const closeHandler = () => {
    modals.setUserEditModal(false);
  }

  return {
    showModal: modals.userEditModal,
    person, canEdit, closeHandler
  }
}