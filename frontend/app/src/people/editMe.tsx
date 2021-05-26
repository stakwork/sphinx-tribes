import React, { useEffect, useState } from "react";
import { useStores } from "../store";
import { useObserver } from "mobx-react-lite";
import {
  EuiModal,
  EuiModalBody,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiOverlayMask,
} from "@elastic/eui";
import Form, { FormField } from "../form";
import ConfirmMe from "./confirmMe";
import type {MeInfo} from '../store/ui'
import api from '../api'

const meSchema: FormField[] = [
  {
    name: "img",
    label: "Image",
    type: "img",
  },
  {
    name: "pubkey",
    label: "Pubkey",
    type: "text",
    readOnly: true,
  },
  {
    name: "owner_alias",
    label: "Name",
    type: "text",
    required: true,
  },
  {
    name: "description",
    label: "Description",
    type: "text",
  },
  {
    name: "price_to_meet",
    label: "Price to Meet",
    type: "number",
  },
  {
    name: "id",
    label: "ID",
    type: "hidden",
  },
];

const host = window.location.host.includes("localhost")
  ? "localhost:5002"
  : window.location.host;

export default function EditMe(props: any) {
  const { ui, main } = useStores();

  const [loading, setLoading] = useState(false);

  function closeModal() {
    ui.setEditMe(false);
    ui.setMeInfo(null)
  }

  async function testChallenge(chal: string) {
    try {
      const me:MeInfo = await api.get(`poll/${chal}`)
      if(me && me.pubkey) {
        ui.setMeInfo(me)
        ui.setEditMe(true)
      }
    } catch(e) {
      console.log(e)
    }
  }

  useEffect(()=>{
    try {
      var urlObject = new URL(window.location.href);
      var params = urlObject.searchParams;
      const chal = params.get('challenge')
      if(chal) {
        testChallenge(chal)
      }
    } catch(e) {}
  })

  async function submitForm(v) {
    console.log(v);
    const info = ui.meInfo as any;
    if (!info) return console.log("no meInfo");
    setLoading(true);
    const URL = info.url.startsWith('http') ? info.url : `https://${info.url}`
    const r = await fetch(URL + "/profile", {
      method: "POST",
      body: JSON.stringify({ 
        ...v, host, 
        price_to_meet: parseInt(v.price_to_meet) 
      }),
      headers: {
        "x-jwt": info.jwt,
        "Content-Type": "application/json"
      },
    });
    if (!r.ok) {
      setLoading(false);
      return alert("Failed to create profile");
    }
    await main.getPeople();
    ui.setEditMe(false);
    ui.setMeInfo(null);
    setLoading(false);
  }
  return useObserver(() => {
    if (!ui.editMe) return <></>;

    let verb = "Create";
    if (ui.meInfo && ui.meInfo.id) verb = "Edit";

    let initialValues = {};
    if (ui.meInfo) {
      initialValues = {
        id: ui.meInfo.id || 0,
        pubkey: ui.meInfo.pubkey,
        owner_alias: ui.meInfo.alias || "",
        img: ui.meInfo.photo_url || "",
        price_to_meet: ui.meInfo.price_to_meet || 0,
        description: ui.meInfo.description || "",
      };
    }

    return (
      <EuiOverlayMask>
        <EuiModal onClose={closeModal} initialFocus="[name=popswitch]">
          <EuiModalHeader>
            <EuiModalHeaderTitle>{`${verb} My Profile`}</EuiModalHeaderTitle>
          </EuiModalHeader>
          <EuiModalBody>
            <div>
              {!ui.meInfo && <ConfirmMe />}
              {ui.meInfo && (
                <Form
                  loading={loading}
                  onSubmit={submitForm}
                  schema={meSchema}
                  initialValues={initialValues}
                />
              )}
            </div>
          </EuiModalBody>
        </EuiModal>
      </EuiOverlayMask>
    );
  });
}
