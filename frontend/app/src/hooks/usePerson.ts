import { useStores } from 'store';
import { Person } from 'store/main';

export const usePerson = (id) => {
  const { main, ui } = useStores();
  const { meInfo } = ui || {};

  let person: Person | undefined;

  if(main.personWanteds.length) {
    const pubkey = main.personWanteds[0].body?.assignee.owner_pubkey;
    person = (main.people || []).find((f) => f.owner_pubkey === pubkey);

    // if(person) person.extras.tickets = main.peopleWanteds;
  } else {
    person = (main.people || []).find((f) => f.id === id);
  }

  const canEdit = meInfo?.id === person?.id;

  return {
    person: canEdit ? meInfo : person,
    canEdit
  };
};
