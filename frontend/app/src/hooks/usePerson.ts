import { useStores } from 'store';
import { Person } from 'store/main';

export const usePerson = (id) => {
  const { main, ui } = useStores();
  const { meInfo } = ui || {};

  let person: Person | undefined;

  if(main.personWanteds.length) {
    const pid = main.personWanteds[0].person.id;
    person = (main.people || []).find((f) => f.id === pid);
    
  } else {
    person = (main.people || []).find((f) => f.id === id);
  }

  const canEdit = meInfo?.id === person?.id;

  return {
    person: canEdit ? meInfo : person,
    canEdit
  };
};
