import React, { useEffect, useState, useRef } from "react";
import { useStores } from "../store";
import { useObserver } from "mobx-react-lite";
import {
  EuiModal,
  EuiModalBody,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiOverlayMask,
} from "@elastic/eui";
import Form from "../form";
import ConfirmMe from "./confirmMe";
import type { MeInfo, MeData } from '../store/ui'
import { emptyMeInfo } from '../store/ui'
import { meSchema } from '../form/schema'
import api from '../api'
import styled, { css } from "styled-components";
import { getHostIncludingDockerHosts } from "../host";

export default function EditMe(props: any) {
  const { ui, main } = useStores();

  const [loading, setLoading] = useState(false);
  const scrollDiv: any = useRef(null)

  function closeModal() {
    ui.setEditMe(false);
    ui.setMeInfo(null);
  }

  async function testChallenge(chal: string) {
    try {
      const me: MeInfo = await api.get(`poll/${chal}`);
      if (me && me.pubkey) {
        ui.setMeInfo(me);
        ui.setEditMe(true);
      }
    } catch (e) {
      console.log(e);
    }
  }

  useEffect(() => {
    try {
      var urlObject = new URL(window.location.href);
      var params = urlObject.searchParams;
      const chal = params.get("challenge");
      if (chal) {
        testChallenge(chal);
      }
    } catch (e) { }
  }, []);

  async function submitForm(body) {
    console.log('SUBMIT FORM', body);
    const info = ui.meInfo as any;
    if (!info) return console.log("no meInfo");
    setLoading(true);
    try {
      const URL = info.url.startsWith("http") ? info.url : `https://${info.url}`;
      const r = await fetch(URL + "/profile", {
        method: "POST",
        body: JSON.stringify({
          // use docker host (tribes.sphinx), because relay will post to it
          host: getHostIncludingDockerHosts(),
          ...body,
          price_to_meet: parseInt(body.price_to_meet),
        }),
        headers: {
          "x-jwt": info.jwt,
          "Content-Type": "application/json",
        },
      });
      if (!r.ok) {
        setLoading(false);
        return alert("Failed to create profile");
      }

      closeModal()
    } catch (e) {
      console.log('e', e)
    }
    setLoading(false);

  }


  return useObserver(() => {
    if (!ui.editMe) return <></>;

    let verb = "Create";
    if (ui.meInfo && ui.meInfo.id) verb = "Edit";

    let initialValues: MeData = emptyMeInfo;

    if (ui.meInfo) {
      initialValues.id = ui.meInfo.id || 0
      initialValues.pubkey = ui.meInfo.pubkey
      initialValues.owner_alias = ui.meInfo.alias || ""
      initialValues.photo_url = ui.meInfo.photo_url || ""
      initialValues.price_to_meet = ui.meInfo.price_to_meet || 0
      initialValues.description = ui.meInfo.description || ""
      initialValues.extras = ui.meInfo.extras || {}

    }

    return (
      <EuiOverlayMask>
        <EuiModal onClose={closeModal}
          style={{
            minWidth: 300,
            minHeight: 460,
            maxWidth: 460,
            maxHeight: 500,
            height: '50vh',
            width: '50vw',
          }}
          initialFocus="[name=popswitch]">
          <EuiModalHeader>
            <EuiModalHeaderTitle>{`${verb} My Profile`}</EuiModalHeaderTitle>
          </EuiModalHeader>
          <EuiModalBody style={{ padding: 0 }}>
            <B ref={scrollDiv}>

              {!ui.meInfo && <ConfirmMe />}
              {ui.meInfo && (
                <Form
                  paged={true}
                  loading={loading}
                  onSubmit={submitForm}
                  scrollDiv={scrollDiv}
                  schema={meSchema}
                  initialValues={initialValues}
                  extraHTML={
                    ui.meInfo.verification_signature
                      ? {
                        twitter: `<span>Post this to your twitter account to verify:</span><br/><strong>Sphinx Verification: ${ui.meInfo.verification_signature}</strong>`,
                      }
                      : {}
                  }
                />
              )}
            </B>
          </EuiModalBody>
        </EuiModal>
      </EuiOverlayMask >
    );
  });
}


const EnvWithScrollBar = ({ thumbColor, trackBackgroundColor }) => css`
  scrollbar-color: ${thumbColor} ${trackBackgroundColor}; // Firefox support
  scrollbar-width: thin;

  &::-webkit-scrollbar {
    width: 6px;
    height: 100%;
  }

  &::-webkit-scrollbar-thumb {
    background-color: ${thumbColor};
    background-clip: content-box;
    border-radius: 5px;
      border: 1px solid ${trackBackgroundColor};
  }

  &::-webkit-scrollbar-corner,
  &::-webkit-scrollbar-track {
    background-color: ${trackBackgroundColor};
  }
}

`

const B = styled.div`
  height:calc(100% - 4px);
  width: calc(100% - 4px);
  overflow-y:auto;
  padding:0 20px;
  box-sizing:border-box;
  ${EnvWithScrollBar({
  thumbColor: '#5a606c',
  trackBackgroundColor: 'rgba(0,0,0,0)',
})}
`

