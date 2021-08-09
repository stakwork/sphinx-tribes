import React, { useRef, useState } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { getHost } from "../host";
import qrCode from "../utils/invoice-qr-code.svg";
import { EuiCheckableCard, EuiButton } from "@elastic/eui";
import { formatPrice } from '../helpers';
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
    <EuiCheckableCard
      className="col-md-6 col-lg-6 ml-2 mb-2"
      id={id + ""}
      label={owner_alias}
      name={owner_alias}
      value={id + ""}
      checked={selected}
      style={{ border: '1px solid #45b9f6' }}
      onChange={() => console.log('change')}
    >
      <Content
        onClick={() => select(id, unique_name)}
        style={{ borderRadius: 5 }}
      >
        <Left>
          <Row className="item-cont" style={{ padding: 10 }}>

            <Img src={img || '/static/sphinx.png'} />

            <Left
              style={{ padding: "0 0 0 20px", maxWidth: "calc(100% - 100px)" }}
            >
              <Row>
                <Title className="tribe-title">{owner_alias}</Title>
              </Row>
              <Description
                style={{ minHeight: 20 }}
              >
                {formatPrice(price_to_meet)}
              </Description>
              {/* <TagsWrap>
              {tags.map((t:string)=> <Tag key={t}>{t}</Tag>)}
            </TagsWrap> */}
            </Left>
          </Row>

          {twitter_confirmed && twitterUsername && (
            <RowWrap>
              <Row
                style={{
                  marginTop: 20,
                  marginBottom: 20,
                  color: "white",
                  fontWeight: "bold",
                  fontSize: 14,
                }}
              >
                <img
                  src="/static/twitter.png"
                  style={{ marginLeft: 28 }}
                  height="32"
                  width="32"
                  alt="twitter"
                />
                <span style={{ marginLeft: 12 }}>{`@${twitterUsername}`}</span>
              </Row>
            </RowWrap>
          )}
          <Intro>{description}</Intro>
        </Left>
      </Content>
    </EuiCheckableCard>
  );
}
interface ContentProps {
  selected: boolean;
}
const Content = styled.div`
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  // max-width: 100%;
  & h3 {
    color: #fff;
  }
  &:hover h3 {
    color: white;
  }
  
    & h5{
      color:#aaa;
    }
    &:hover h5{
      color:#bebebe;
    }
  }
  & button img {
  }
`;
const QRWrapWrap = styled.div`
  display: flex;
  justify-content: center;
`;
const QRWrap = styled.div`
  background: white;
  padding: 5px;
`;
const Left = styled.div`
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
const RowWrap = styled.div``;
const Title = styled.h3`
  color:#fff;
  margin-right: 12px;
  font-size: 22px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  min-height: 24px;
`;
const Description = styled.div`
  
  font-size: 12px;
  font-weight: 400;
  color: #6b7a8d;
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
  height: 70px;
  width: 70px;
  border-radius: 50%;
  position: relative;
`;
const Tokens = styled.div`
  display: flex;
  align-items: center;
`;
const TagsWrap = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  margin-top: 10px;
`;
const Tag = styled.h5`
  margin-right: 10px;
`;
const Intro = styled.div`
  color: white;
  font-size: 14px;
  margin: 5px;
  padding: 10px;
  max-width: 400px;
  max-height:100px;
  overflow:auto;
  background: #ffffff21;
  border-radius: 5px;
`;
