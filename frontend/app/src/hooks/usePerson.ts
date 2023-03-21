import { useStores } from 'store';
import { Person } from 'store/main';

export const usePerson = (id) => {
  const { main, ui } = useStores();
  const { meInfo } = ui || {};

  const person: Person | undefined = (main.people || []).find((f) => f.id === id);

  const canEdit = meInfo?.id === person?.id;
  return {
    person: canEdit ? meInfo : person,
    canEdit
  };
};
