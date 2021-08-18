import React, { useState, useEffect } from "react";
import { QRCode } from "react-qr-svg";
import { EuiLoadingSpinner } from "@elastic/eui";
import styled from "styled-components";
import api from "../api";
import { useStores } from "../store";
import type { MeInfo } from "../store/ui";
import { getHost } from "../host";

const host = getHost();
function makeQR(challenge: string, ts: string) {
  return `sphinx.chat://?action=auth&host=${host}&challenge=${challenge}&ts=${ts}`;
}

export default function ConfirmMe() {
  const { ui } = useStores();
  const [challenge, setChallenge] = useState("");
  const [ts, setTS] = useState("");

  async function startPolling(challenge: string) {
    let ok = true;
    let i = 0;
    while (ok) {
      await sleep(3000);
      try {
        const me: MeInfo = await api.get(`poll/${challenge}`);
        console.log(me);
        if (me && me.pubkey) {
          ui.setMeInfo(me);
          setChallenge("");
          ok = false;
          break;
        }
        i++;
        if (i > 100) ok = false;
      } catch (e) { }
    }
  }
  async function getChallenge() {
    const res = await api.get("ask");
    if (res.challenge) {
      setChallenge(res.challenge);
      startPolling(res.challenge);
    }
    if (res.ts) {
      setTS(res.ts);
    }
  }
  useEffect(() => {
    getChallenge();
  }, []);

  const qrString = makeQR(challenge, ts);
  return (
    <ConfirmWrap>
      {!challenge && <EuiLoadingSpinner size="xl" style={{ marginTop: 60 }} />}
      {challenge && (
        <InnerWrap>
          <P>Scan QR or click to open Sphinx</P>
          <QrWrap>
            <QRCode
              bgColor="#FFFFFF"
              fgColor="#000000"
              level="Q"
              style={{ width: 209 }}
              value={qrString}
            />
          </QrWrap>
          <LinkWrap>
            <a href={qrString} className="btn join-btn">
              <img
                style={{ width: 13, height: 13, marginRight: 8 }}
                src="/static/launch-24px.svg"
                alt=""
              />
              Open Sphinx
            </a>
          </LinkWrap>
        </InnerWrap>
      )}
    </ConfirmWrap>
  );
}

const ConfirmWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
  min-height: 250px;
`;
const InnerWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
`;
const LinkWrap = styled.div`
  width: 100%;
  text-align: center;
  margin: 20px 0;
  & a {
    width: 115px;
    position: relative;
    margin-left: 25px;
  }
`;
const P = styled.p`
  color: white;
`;
const QrWrap = styled.div`
  padding: 8px;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
`;

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
