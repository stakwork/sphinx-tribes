import React, { useRef, useState } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { getHost } from "../host";
import qrCode from "../utils/invoice-qr-code.svg";
import { EuiCheckableCard, EuiButton } from "@elastic/eui";

const host = getHost();
function makeQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export default function Person({
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
}: any) {
  const [showQR, setShowQR] = useState(false);
  const qrString = makeQR(owner_pubkey);

  const twitterUsername = (extras && extras.twitter) || null;

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
  return (
    <EuiCheckableCard
      className="col-md-6 col-lg-6 ml-2 mb-2"
      id={id + ""}
      label={owner_alias}
      name={owner_alias}
      value={id + ""}
      checked={selected}
      onChange={() => select(id, unique_name)}
    >
      <Content
        onClick={() => select(selected ? "" : id, unique_name)}
        style={{
          height: selected ? "auto" : 100,
        }}
        selected={selected}
      >
        <Left>
          <Row className="item-cont">
            {img ? (
              <Img src={img} />
            ) : (
              <div className="placeholder-img-tribe"></div>
            )}
            <Left
              style={{ padding: "0 0 0 20px", maxWidth: "calc(100% - 100px)" }}
            >
              <Row
                style={
                  selected
                    ? { flexDirection: "column", alignItems: "flex-start" }
                    : {}
                }
              >
                <Title className="tribe-title">{owner_alias}</Title>
              </Row>
              <Description
                oneLine={selected ? false : true}
                style={{ minHeight: 20 }}
              >
                {description}
              </Description>
              {/* <TagsWrap>
              {tags.map((t:string)=> <Tag key={t}>{t}</Tag>)}
            </TagsWrap> */}
            </Left>
          </Row>
          {price_to_meet && (
            <RowWrap>
              <Row
                style={{
                  marginTop: 20,
                  marginBottom: 20,
                  justifyContent: "space-evenly",
                  color: "white",
                  fontWeight: "bold",
                }}
              >
                {`Price to Meet: ${price_to_meet} sats`}
              </Row>
            </RowWrap>
          )}
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
          <RowWrap>
            <Row style={{ marginBottom: 20, justifyContent: "space-evenly" }}>
              <a href={qrString}>
                <EuiButton
                  onClick={add}
                  fill={true}
                  style={{
                    backgroundColor: "#6089ff",
                    borderColor: "#6089ff",
                    color: "white",
                    fontWeight: 600,
                    fontSize: 12,
                  }}
                  aria-label="add"
                >
                  ADD
                </EuiButton>
              </a>
              <EuiButton
                onClick={toggleQR}
                style={{
                  borderColor: "#6B7A8D",
                  color: "white",
                  fontWeight: 600,
                  fontSize: 12,
                }}
                iconType={qrCode}
                aria-label="qr-code"
              >
                QR CODE
              </EuiButton>
            </Row>
          </RowWrap>
          {showQR && (
            <QRWrapWrap>
              <QRWrap className="qr-wrap float-r">
                <QRCode
                  bgColor={"#FFFFFF"}
                  fgColor="#000000"
                  level="Q"
                  style={{ width: 209 }}
                  value={qrString}
                />
              </QRWrap>
            </QRWrapWrap>
          )}
          <Intro>{description}</Intro>
          {/* <Row style={{ color:"#6B7A8D", fontSize:12, fontWeight:"bold", padding:"10px 10px 0px 10px" }}>
          INTERESTS
        </Row>
        <Row style={{color:"white", fontSize:14, margin:"0px 10px 10px 10px"}}>
          {tagsString}
        </Row> */}
          <div
            className="expand-part"
            style={selected ? { opacity: 1 } : { opacity: 0 }}
          >
            <div className="colapse-button">
              <img src="/static/keyboard_arrow_up-black-18dp.svg" alt="" />
            </div>
          </div>
        </Left>
      </Content>
    </EuiCheckableCard>
  );
}
interface ContentProps {
  selected: boolean;
}
const Content = styled.div<ContentProps>`
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  max-width: 100%;
  & h3 {
    color: #fff;
  }
  &:hover h3 {
    color: white;
  }
  ${(p) =>
    p.selected
      ? `
    & h5{
      color:#cacaca;
    }
  `
      : `
    & h5{
      color:#aaa;
    }
    &:hover h5{
      color:#bebebe;
    }
  `}
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
  margin-right: 12px;
  font-size: 22px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  min-height: 24px;
`;
interface DescriptionProps {
  oneLine: boolean;
}
const Description = styled.div<DescriptionProps>`
  font-weight: normal;
  line-height: 20px;
  font-size: 15px;
  font-weight: 500;
  color: #6b7a8d;
  ${(p) =>
    p.oneLine &&
    `
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow:hidden;
  `}
`;
interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url("${(p) => p.src}");
  background-position: center;
  background-size: cover;
  height: 90px;
  width: 90px;
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
  margin: 10px;
`;
