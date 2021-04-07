import React, {useState} from 'react'
import { useStores } from "../store";
import { useObserver } from "mobx-react-lite";
import {
  EuiModal,
  EuiModalBody,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiOverlayMask,
} from '@elastic/eui';
import Form, {FormField} from '../form'
import ConfirmMe from './confirmMe'

const meSchema: FormField[] = [
  {
    name:'owner_alias',
    label:'Name',
    type:'text'
  },
  {
    name:'description',
    label:'Description',
    type:'text'
  }
]

export default function EditMe(props:any) {
  const { ui } = useStores();
  const [me, setMe] = useState({})

  function closeModal(){
    ui.setEditMe(false)
  }
  return useObserver(() => {
    if(!ui.editMe) return <></>
    return <EuiOverlayMask>
      <EuiModal onClose={closeModal} initialFocus="[name=popswitch]">
        <EuiModalHeader>
          <EuiModalHeaderTitle>My Profile</EuiModalHeaderTitle>
        </EuiModalHeader>
        <EuiModalBody>
          {!me && <ConfirmMe setMe={setMe} />}
          {me && <Form 
            schema={meSchema}
          />}
        </EuiModalBody>
      </EuiModal>
    </EuiOverlayMask>
  });
}
