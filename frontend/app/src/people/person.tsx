import React, { useRef, useState } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { getHost } from "../host";
import qrCode from "../utils/invoice-qr-code.svg";
import { EuiCheckableCard, EuiButton } from "@elastic/eui";
import { formatPrice } from '../helpers';
import { colors } from "../colors";
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
    twitter_confirmed,
  } = props

  const [showQR, setShowQR] = useState(false);
  const qrString = makeQR(owner_pubkey);
  const c = colors['light']

  const twitterUsername = (extras && extras.twitter && extras.twitter.handle) || (extras && extras.twitter) || null;

  let tagsString = "";
  tags.forEach((t: string, i: number) => {
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


  // return <div style={{ color: '#fff' }}>{owner_alias}</div>
  return (
    <Wrap onClick={() => select(id, unique_name)}>
      <EuiCheckableCard
        id={id + ""}
        label={owner_alias}
        name={owner_alias}
        value={id + ""}
        checked={selected}
        onChange={() => console.log('change')}
        style={{ background: c.background, color: c.text1 }}
      >
        <Content>

          <Row className="item-cont" style={{ padding: 10 }}>

            <Img src={img || '/static/sphinx.png'} />

            <Left>
              <Row>
                <Title>{owner_alias}</Title>
              </Row>
              <Row>
                <Description>
                  {description}
                </Description>
              </Row>
            </Left>
          </Row>

        </Content>
      </EuiCheckableCard>
    </Wrap>
  );
}
interface ContentProps {
  selected: boolean;
}
const Content = styled.div`
  display: flex;
  align-items:center;
  width:250px;
`;
const Left = styled.div`
  height: 100%;
  max-width: 100%;
  margin-left:20px;
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const Wrap = styled.div`
  cursor:pointer;
`;

const B = styled.div`
  height: 100%;
  max-width: 100%;
  display: flex;
  flex-direction: column;
  flex: 1;
`;
const Row = styled.div`
  display: flex;
  align-items: center;
`;
const Title = styled.h3`
  font-size: 22px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  font-weight: 500;
  font-size: 17px;
  color: #3C3F41;
`;
const Description = styled.div`
  font-size: 12px;
  font-weight: 400;
  line-height:19px;
  color: #8E969C;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow:hidden;
  
`;
interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url("${(p) => p.src}");
  background-position: center;
  background-size: cover;
  height: 80px;
  width: 80px;
  border-radius: 50%;
  position: relative;
`;