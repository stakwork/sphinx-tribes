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
    name:'img',
    label:'Image',
    type:'img'
  },
  {
    name:'pubkey',
    label:'Pubkey',
    type:'text',
    readOnly:true
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
  },
  {
    name:'price_to_meet',
    label:'Price to Meet',
    type:'number'
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
        <EuiModalBody >
          <div>
            {!ui.meInfo && <ConfirmMe />}
            {ui.meInfo && <Form 
              schema={meSchema}
              initialValues={{
                pubkey: ui.meInfo.pubkey,
                owner_alias: ui.meInfo.alias,
                img: ui.meInfo.photoUrl
              }}
            />}
          </div>
        </EuiModalBody>
      </EuiModal>
    </EuiOverlayMask>
  });
}
