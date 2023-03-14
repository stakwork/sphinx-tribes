import { getHost } from "config";
import { useMemo } from "react";
import { useHistory } from "react-router-dom";
import { useStores } from "store";

//TODO: mv into utils
const host = getHost();
function makeQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export const useUserInfo = () => {
  const { main, ui } = useStores();
  const history = useHistory();
  const { meInfo } = ui || {};
  const personId = ui.selectedPerson;

  const person: any = main.people && main.people.length && main.people.find((f) => f.id === personId);

  const { id, img, owner_alias, extras, owner_pubkey } = person || {};
  const canEdit = id === meInfo?.id;
  function goBack() {
    ui.setSelectingPerson(0);
    history.goBack();
  }
  const qrString = makeQR(owner_pubkey || '');

  const defaultPic = '/static/person_placeholder.png';
  const userImg = img || defaultPic;
  
  function logout() {
    ui.setEditMe(false);
    ui.setMeInfo(null);
    main.getPeople({ resetPage: true });
    goBack();
  }

  return {
    canEdit, 
    goBack, 
    userImg, 
    owner_alias, 
    logout, 
    person,
    qrString
  }
}