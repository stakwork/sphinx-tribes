import React, { useRef, useState } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { EuiCard } from "@elastic/eui";
import Tag from "./tag";
import moment from "moment";
import { getHostIncludingDockerHosts } from "../host";

function makeQR(uuid: string) {
  return `sphinx.chat://?action=tribe&uuid=${uuid}&host=${getHostIncludingDockerHosts()}`;
}

export default function Tribe({
  uuid,
  name,
  img,
  tags,
  description,
  selected,
  select,
  created,
  owner_alias,
  price_to_join,
  price_per_message,
  member_count,
  last_active,
  unique_name,
  preview,
}: any) {
  const showTags = tags && tags.length && tags.length > 0 ? true : false;
  const qrString = makeQR(uuid);
  const [copied, setCopied] = useState(false);

  const lastActiveM = last_active
    ? moment(last_active * 1000)
    : moment().subtract(1, "months");
  const lastActive = lastActiveM.format("MMM D HH:mm");

  function copyString(e, text: any) {
    e.stopPropagation();
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2700);
  }

  return (
    <EuiCard
					style={{backgroundColor: "#1a2430"}}
      className="col-md-6 col-lg-6  mb-2"
      id={uuid}
      title={""}
    >
      <Content
        onClick={() => select(selected ? "" : uuid, unique_name)}
        style={{
          height: selected ? "auto" : 100,
        }}
        selected={selected}
      >
        <Left>
          <Row className="item-cont">
            <Img src={img || "/static/placeholder.svg"} />
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
                <Title className="tribe-title">{name}</Title>
              </Row>
              <Description
                oneLine={!selected }
                style={{ minHeight: 20 }}
              >
                {description}
              </Description>

              {showTags && (
                <Tokens className="tags-wrapper">
                  {tags.map((t: string) => (
                    <Tag type={t} key={t} />
                  ))}
                  {tags.length > 3 && (
                    <span className="more-tags-available">more...</span>
                  )}
                </Tokens>
              )}
            </Left>
          </Row>
          <div
            className="expand-part"
            style={selected ? { opacity: 1 } : { opacity: 0 }}
          >
            <div className="section-separator"></div>
            <div className="row info-section">
              <div className="col-4 col-sm-4 col-md-4 col-lg-4">
                <div className="uppercase">Last Activity:</div>
                <div className="lighter-color">{lastActive}</div>
              </div>
              <div className="col-4 col-sm-4 col-md-4 col-lg-4 text-center">
                {(member_count ? true : false) && (
                  <>
                    <div className="uppercase">Members:</div>
                    <div className="lighter-color">{member_count}</div>
                  </>
                )}
              </div>
              <div className="col-4 col-sm-4 col-md-4 col-lg-4 text-right">
                <div className="uppercase">Admin:</div>
                <div className="lighter-color">{owner_alias}</div>
              </div>
            </div>
            <div className="section-separator"></div>

            <div className="row">
              <div className="col-4 col-sm-4 col-md-4 col-lg-4 qr-left">
                <div className="text-right">
                  <img
                    style={{ width: 100 }}
                    src="/static/scan_notification.svg"
                    alt=""
                  />
                </div>
                <div className="text-right info">
                  <div>Price to join</div>
                  <div className="lighter-color">{price_to_join || 0}</div>
                </div>
                <div className="section-separator"></div>
                <div className="text-right info">
                  <div>Price per msg</div>
                  <div className="lighter-color">{price_per_message || 0}</div>
                </div>
                <div className="section-separator"></div>
                <a
                  onClick={(e) => e.stopPropagation()}
                  href={qrString}
                  className="btn join-btn"
                >
                  <img
                    style={{ width: 13, height: 13, marginRight: 8 }}
                    src="/static/launch-24px.svg"
                    alt=""
                  />
                  Join
                </a>
              </div>
              <div className="col-8 col-sm-8 col-md-8 col-lg-8">
                {selected && (
                  <QRWrap className="qr-wrap float-r">
                    <QRCode
                      bgColor={selected ? "#FFFFFF" : "#666"}
                      fgColor="#000000"
                      level="Q"
                      style={{ width: 209 }}
                      value={qrString}
                    />
                  </QRWrap>
                )}
                <div className="below-qr">
                  <a
                    onClick={(e) => e.stopPropagation()}
                    href={`https://${preview}?tribe=${uuid}`}
                    target={"_blank"}
                    className={`preview-btn ${!preview ? "btn-disabled" : ""}`}
                  >
                    Preview
                  </a>
                  <button
                    className="copy-btn"
                    onClick={(e) => copyString(e, qrString)}
                  >
                    {copied ? "Copied!" : "Copy Link"}
                  </button>
                </div>
              </div>
            </div>

            <div className="colapse-button">
              <img src="/static/keyboard_arrow_up-black-18dp.svg" alt="" />
            </div>
          </div>
        </Left>
      </Content>
    </EuiCard>
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
const Description = styled.h5<DescriptionProps>`
  font-weight: normal;
  line-height: 20px;
	align-self: start;
	font-size: 10px;
	text-align: left;
  ${(p) =>
    p.oneLine &&
    `
		max-width: 100%;
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
  border-radius: 5px;
  position: relative;
`;

const Tokens = styled.div`
  display: flex !important;
	align-content: start;
`;
