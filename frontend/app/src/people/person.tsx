import React, { useRef, useState } from "react";
import { QRCode } from "react-qr-svg";
import styled from "styled-components";
import { getHost } from "../host";
import { useObserver } from 'mobx-react-lite'
import { colors } from "../colors";
import { Button, Divider, Modal } from '../sphinxUI/index'
const host = getHost();
function makeQR(pubkey: string) {
  return `sphinx.chat://?action=person&host=${host}&pubkey=${pubkey}`;
}

export default function Person(props: any) {

  const {
    hideActions,
    small,
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

  const c = colors['light']

  let tagsString = "";
  tags && tags.forEach((t: string, i: number) => {
    if (i !== 0) tagsString += ",";
    tagsString += t;
  });


  return useObserver(() => {


    const qrString = makeQR(owner_pubkey);

    function renderPersonCard() {
      if (small) {
        return <Wrap onClick={() => select(id, unique_name)} style={{
          background: selected ? '#F2F3F5' : '#fff',

        }}>
          <div>
            <Img src={img || '/static/sphinx.png'} style={hideActions && { width: 56, height: 56 }} />
          </div>
          <R>
            <Title style={hideActions && { fontSize: 17 }}>{owner_alias}</Title>
            <Description>
              {description}
            </Description>
            {!hideActions &&
              <Row style={{ justifyContent: 'space-between', alignItems: 'center' }}>
                <div></div>
                {(!hideActions && owner_pubkey) ?
                  <>
                    <a href={qrString}>
                      <Button
                        text='Connect'
                        color='white'
                        leadingIcon={'open_in_new'}
                        iconSize={16}
                        onClick={(e) => e.stopPropagation()}
                      />
                    </a>
                  </> : <div style={{ height: 30 }} />
                }
              </Row>
            }
            <Divider style={{ marginTop: 20 }} />
          </R>
        </Wrap>
      }
      // desktop mode
      return <DWrap onClick={() => select(id, unique_name)}>
        <div>
          <Img style={{ height: 210, width: '100%', borderRadius: 0 }} src={img || '/static/sphinx.png'} />
          <div style={{ padding: 10 }}>
            <DTitle>{owner_alias}</DTitle>
            <DDescription>
              {description}
            </DDescription>
          </div>
        </div>
        <div>
          <Divider />
          <Row style={{ justifyContent: 'space-between', alignItems: 'center', height: 50 }}>
            <div />
            {owner_pubkey ?
              <>
                <Button
                  text='Connect'
                  color='clear'
                  endingIcon={'open_in_new'}
                  iconSize={16}
                  onClick={(e) => {
                    setShowQR(true)
                    e.stopPropagation()
                  }}
                />
              </> : <div />
            }
          </Row>
        </div>
      </DWrap>
    }

    return (
      <>
        {renderPersonCard()}

        <Modal
          visible={showQR}
          close={() => setShowQR(false)}
        >
          <InnerWrap>
            <QRCode
              bgColor={"#FFFFFF"}
              fgColor="#000000"
              // level="Q"
              style={{ width: 200, height: 200 }}
              value={qrString}
            />
            <div style={{ marginTop: 10 }}>Scan with your Sphinx Mobile App</div>
          </InnerWrap>
        </Modal>
      </>
    );
  })
}
const InnerWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  width: 100%;
  padding:50px 0 30px;
`;
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
const DWrap = styled.div`
        cursor:pointer;
        height:350px;
        width:210px;
        display:flex;
        flex-direction:column;
        justify-content:space-between;
        background:#fff;
        margin-bottom:20px;
        box-shadow: 0px 1px 2px rgba(0, 0, 0, 0.15);
        border-radius: 4px;
        `;

const R = styled.div`
        width:calc(100% - 80px);
        margin-left:20px;
        `;


const Row = styled.div`
        display: flex;
        width:100%;
        `;

const Title = styled.h3`
        font-weight: 500;
        font-size: 20px;
        line-height: 19px;
        /* or 95% */


        /* Text 2 */

        color: #3C3F41;
        `;

const DTitle = styled.h3`
font-weight: 500;
font-size: 17px;
line-height: 19px;

        color: #3C3F41;
        `;
const Description = styled.div`
        font-size: 15px;
        color: #5F6368;
        white-space: nowrap;
        height:26px;
        text-overflow: ellipsis;
        overflow:hidden;
        margin-bottom:10px;
        `;

const DDescription = styled.div`
        font-size: 12px;
        line-height: 18px;
        color: #5F6368;
        white-space: nowrap;
        // height:26px;
        text-overflow: ellipsis;
        overflow:hidden;
        margin-bottom:10px;
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
          border-radius: 50%;
          position: relative;
          `;