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
  },
  {
    name:'id',
    label:'ID',
    type:'hidden'
  }
]

const host = window.location.host.includes('localhost')?'localhost:5002':window.location.host

export default function EditMe(props:any) {
  const { ui, main } = useStores();

  const [loading,setLoading] = useState(false)

  function closeModal(){
    ui.setEditMe(false)
  }

  async function submitForm(v) {
    console.log(v)
    const info = ui.meInfo as any
    if(!info) return console.log("no meInfo")
    setLoading(true)
    const url = info.url
    const jwt = info.jwt
    const r = await fetch(url+'/profile', {
      method:'POST',
      body:JSON.stringify({...v, host}),
      headers:{
        'x-jwt': jwt
      }
    })
    if(!r.ok) {
      setLoading(false)
      return alert('Failed to create profile')
    }
    await main.getPeople()
    ui.setEditMe(false)
    ui.setMeInfo(null)
    setLoading(false)
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
              loading={loading}
              onSubmit={submitForm}
              schema={meSchema}
              initialValues={{
                id: ui.meInfo.id || 0,
                pubkey: ui.meInfo.pubkey,
                owner_alias: ui.meInfo.alias,
                img: ui.meInfo.photo_url,
              }}
            />}
          </div>
        </EuiModalBody>
      </EuiModal>
    </EuiOverlayMask>
  });
}
