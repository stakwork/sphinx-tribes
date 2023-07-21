import { useStores } from 'store';
import { Person } from 'store/main';

export const usePerson = (id: any) => {
  const { main, ui } = useStores();
  const { meInfo } = ui || {};

  let person: Person | undefined;

  if (main.personAssignedWanteds.length) {
    const pubkey = main.personAssignedWanteds[0].body.assignee.owner_pubkey;
    person = (main.people || []).find((f: any) => f.owner_pubkey === pubkey);
  } else {
    person = (main.people || []).find((f: any) => f.id === id);
  }

  const canEdit = meInfo?.id === person?.id;

  return {
    person: canEdit ? meInfo : person,
    canEdit
  };
};
