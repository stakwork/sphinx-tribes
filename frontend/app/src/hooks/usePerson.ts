import { useStores } from "store";

export const usePerson = (id) => {
  const { main, ui } = useStores();
  const { meInfo } = ui || {};
  // const history = useHistory();
  // const personId = ui.selectedPerson;

  const person: any =
    main.people && main.people.length && main.people.find((f) => f.id === id);

    const canEdit = meInfo?.id === person.id
    return {
      person, 
      canEdit 
    }

}