import React from 'react';
import styled from 'styled-components';
import { Button, Modal } from '../../sphinxUI';
import QR from './QR';
import QrBar from './QrBar';
import { makeConnectQR } from '../../helpers';

export default function ConnectCard(props) {
  const { visible } = props;
  const { person } = props;

  const qrString = makeConnectQR(person?.owner_pubkey);

  return (
    <div onClick={(e)=> e.stopPropagation()} >
    <Modal
      style={props.modalStyle}
      // close={(e) => {
        //   e.stopPropagation();
        //   props.dismiss();
        // }}
        overlayClick={(e: React.SyntheticEvent)=>{
          props.dismiss();
        }}
        visible={visible}>
      {/* <div style={{position:'relative'}}> */}
      <div style={{ textAlign: 'center', paddingTop: 59, width: 310 }}>
        <ImgWrap>
          <W>
            <Icon src={person?.img || '/static/person_placeholder.png'} />
          </W>
        </ImgWrap>
        <div style={{ textAlign: 'center', width: '100%', overflow: 'hidden', padding: '0 50px' }}>
          <N>Discuss this bounty with</N>
          <D>
            <B>{person?.owner_alias} </B>
          </D>

          <QR value={qrString} size={210} type={'connect'} />

          <QrBar value={person?.owner_pubkey} simple style={{ marginTop: 11 }} />

          <a href={qrString}>
            <Button
              text={'Connect with Sphinx'}
              color={'primary'}
              style={{ paddingLeft: 25, margin: '12px 0 40px' }}
              img={'sphinx_white.png'}
              imgSize={27}
              height={48}
              width={'100%'}
              />
          </a>
        </div>
      </div>
        <div style={{position:'absolute' , bottom:'-36px',  width: 310, backgroundColor:'transparent', display:'flex', justifyContent:'center'}} >
          <img src="/static/scan_qr.svg" alt="scan" />
          <div style={{marginLeft:'12px', color:'#FFFFFF'}} >Scan or paste in Sphinx</div>
        </div>
     {/* </div> */}
    </Modal>
   </div>
  );
}

const ImgWrap = styled.div`
  position: absolute;
  top: -45px;
  left: 0px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  pointer-events: none;
  user-select: none;
`;
const B = styled.span`
  font-weight: bold;
`;
const W = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  width: 90px;
  height: 90px;
  border-radius: 80px;
`;
const N = styled.div`
  font-family: Barlow;
  font-style: normal;
  font-weight: 500;
  font-size: 17px;
  line-height: 26px;
  /* or 95% */

  text-align: center;

  /* Text 2 */

  color: #8E969C;
`;

const D = styled.div`
  font-family: Barlow;
  font-style: normal;
  font-size: 20px;
  line-height: 26px;
  /* or 129% */

  text-align: center;

  /* Main bottom icons */

  color: #3c3f41;
  margin-bottom: 20px;
`;

interface IconProps {
  src: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p) => `url(${p.src})`};
  width: 80px;
  height: 80px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: cover; /* Resize the background image to cover the entire container */
  border-radius: 80px;
  overflow: hidden;
`;
