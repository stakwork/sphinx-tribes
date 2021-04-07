import React from 'react'
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
    name:'pubkey',
    label:'Pubkey',
    type:'text',
    readOnly:true
  },
  {
    name:'img',
    label:'Image',
    type:'img'
  },
  {
    name:'owner_alias',
    label:'Name',
    type:'text',
    required: true
  },
  {
    name:'description',
    label:'Description',
    type:'text'
  }
]

export default function EditMe(props:any) {
  const { ui } = useStores();

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
          {!ui.tokens && <ConfirmMe />}
          {ui.tokens && <Form 
            schema={meSchema}
            initialValues={{pubkey:ui.tokens.pubkey}}
          />}
        </EuiModalBody>
      </EuiModal>
    </EuiOverlayMask>
  });
}
