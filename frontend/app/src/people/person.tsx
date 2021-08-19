import React, { useRef, useState } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { getHost } from "../host";
import { useObserver } from 'mobx-react-lite'
import qrCode from "../utils/invoice-qr-code.svg";
import { EuiCheckableCard, EuiButton } from "@elastic/eui";
import { formatPrice } from '../helpers';
import { colors } from "../colors";
import { useIsMobile, useScreenWidth } from "../hooks";
import { Button, Divider } from '../sphinxUI/index'
const host = getHost();
function makeQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export default function Person(props: any) {

  const {
    id,
    img,
    tags,
    description,
    selected,
    select,
    created,
    owner_alias,
    owner_pubkey,
    unique_name,
    price_to_meet,
    extras,
    twitter_confirmed
  } = props

  const [showQR, setShowQR] = useState(false);
  const qrString = makeQR(owner_pubkey);
  const c = colors['light']

  const isMobile = useIsMobile()
  const screenWidth = useScreenWidth()

  const twitterUsername = (extras && extras.twitter && extras.twitter.handle) || (extras && extras.twitter) || null;

  let tagsString = "";
  tags && tags.forEach((t: string, i: number) => {
    if (i !== 0) tagsString += ",";
    tagsString += t;
  });

  function add(e) {
    e.stopPropagation();
  }
  function toggleQR(e) {
    e.stopPropagation();
    setShowQR((current) => !current);
  }

  return useObserver(() => {


    // return <div style={{ color: '#fff' }}>{owner_alias}</div>
    return (
      <Wrap onClick={() => select(id, unique_name)}>
        <div>
          <Img src={img || '/static/sphinx.png'} />
        </div>

        <R>
          <Title>{owner_alias}</Title>
          <Description>
            {description}
          </Description>
          <Row style={{ justifyContent: 'space-between' }}>
            <div>3h ago</div>
            <Button
              text='Connect'
              color='white'

            />
          </Row>
          <Divider style={{ marginTop: 20 }} />
        </R>
      </Wrap>
    );
  })
}
interface ContentProps {
  selected: boolean;
}
const Content = styled.div`
        // display:flex;
        `;
const Wrap = styled.div`
        cursor:pointer;
        padding: 25px;
        padding-bottom:0px;
        display:flex;
        width:100%;
        `;
const R = styled.div`
        width:67%;
        margin-left:20px;
        `;


const Row = styled.div`
        display: flex;
        width:100%;
        `;

const Title = styled.h3`
        font-style: normal;
        font-weight: 500;
        font-size: 20px;
        line-height: 19px;
        /* or 95% */


        /* Text 2 */

        color: #3C3F41;
        `;
const Description = styled.div`
        font-size: 15px;
        color: #5F6368;
        white-space: nowrap;
        height:26px;
        text-overflow: ellipsis;
        overflow:hidden;
        // width:280px;
        `;
interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
          background-image: url("${(p) => p.src}");
          background-position: center;
          background-size: cover;
          height: 96px;
          width: 96px;
          min-width: 96px;
          border-radius: 50%;
          position: relative;
          `;