import { usePerson } from 'hooks';
import { useStores } from 'store';

export const useUserEdit = () => {
  const { ui, modals } = useStores();
  const personId = ui.selectedPerson;
  // const personId = useParams<{personId: string}>();

  const { canEdit, person } = usePerson(Number(personId));

  const closeHandler = () => {
    modals.setUserEditModal(false);
  };

  return {
    modals,
    person,
    canEdit,
    closeHandler
  };
};
