import { EuiModal, EuiOverlayMask } from '@elastic/eui';
import { useState } from 'react';
import React from 'react';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { StartUpModalProps } from 'people/interfaces';
import IconButton from '../../components/common/IconButton2';
import { useStores } from '../../store';
import api from '../../api';
import QR from './QR';

const ModalContainer = styled.div`
  max-height: 274px;
  overflow-y: auto;
  display: flex;
  justify-content: center;
  margin-top: 20px;
  margin-bottom: 50px;
`;

const ButtonContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const QrContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 0;
  overflow-y: hidden;
`;

const QRText = styled.p`
  padding: 0px;
  margin-top: 15px;
  font-size: 0.9rem;
  font-weight: bold;
  width: 60%;
  text-align: center;
`;

const DirectionWrap = styled.div`
  padding: 0px;
  display: flex;
  width: 100%;
`;

const AndroidIosButtonConatiner = styled.div`
  padding: 0px;
  display: flex;
  width: 100%;
  margin-right: 20px;
  justify-content: space-between;
`;
const StartUpModal = ({ closeModal, dataObject, buttonColor }: StartUpModalProps) => {
  const { ui, main } = useStores();
  const [step, setStep] = useState(1);
  const [connection_string, setConnectionString] = useState('');

  async function getConnectionCode() {
    if (!connection_string) {
      const code = await api.get('connectioncodes');
      if (code.connection_string) {
        setConnectionString(code.connection_string);
        main.getPeople({ resetPage: true });
      }
    }
  }

  const StepOne = () => (
    <>
      <ModalContainer>
        <img
          src={
            dataObject === 'getWork'
              ? '/static/create_profile_blue.gif'
              : '/static/create_profile_green.gif'
          }
          height={'274px'}
          alt=""
        />
      </ModalContainer>
      <ButtonContainer>
        <DirectionWrap style={{ justifyContent: 'space-around' }}>
          <IconButton
            text={'I have Sphinx'}
            width={150}
            height={48}
            style={{ marginTop: '20px' }}
            onClick={(e: any) => {
              e.stopPropagation();
              closeModal();
              ui.setShowSignIn(true);
            }}
            textStyle={{
              fontSize: '15px',
              fontWeight: '500'
            }}
            iconStyle={{
              top: '14px'
            }}
            color={buttonColor}
          />

          <IconButton
            text={'Get Sphinx'}
            width={150}
            height={48}
            style={{ marginTop: '20px', textDecoration: 'none' }}
            onClick={(e: any) => {
              e.stopPropagation();
              setStep(step + 1);
            }}
            textStyle={{
              fontSize: '15px',
              fontWeight: '500'
            }}
            iconStyle={{
              top: '14px'
            }}
            color={buttonColor}
          />
        </DirectionWrap>
      </ButtonContainer>
    </>
  );

  const StepTwo = () => (
    <>
      <ModalContainer>
        {connection_string ? (
          <QrContainer>
            <QR size={200} value={ui.connection_string} />
            <QRText>Install the Sphinx app on your phone and then scan this QRcode</QRText>
          </QrContainer>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column' }}>
            <p style={{ textAlign: 'center', fontWeight: 600 }}>Step 1</p>
            <p style={{ textAlign: 'center' }}>Download App</p>
            <AndroidIosButtonConatiner>
              <IconButton
                text={'Android'}
                width={100}
                height={48}
                style={{ marginTop: '20px', textDecoration: 'none' }}
                onClick={() =>
                  window.open('https://play.google.com/store/apps/details?id=chat.sphinx', '_blank')
                }
                textStyle={{
                  fontSize: '15px',
                  fontWeight: '500'
                }}
                iconStyle={{
                  top: '14px'
                }}
                color={buttonColor}
              />
              <IconButton
                text={'IOS'}
                width={100}
                height={48}
                style={{ marginTop: '20px', textDecoration: 'none' }}
                onClick={() => window.open('https://testflight.apple.com/join/QoaCkJn6', '_blank')}
                textStyle={{
                  fontSize: '15px',
                  fontWeight: '500'
                }}
                iconStyle={{
                  top: '14px'
                }}
                color={buttonColor}
              />
            </AndroidIosButtonConatiner>
          </div>
        )}
      </ModalContainer>
      <p style={{ textAlign: 'center', fontWeight: 600 }}>Step 2</p>
      <p style={{ textAlign: 'center' }}>Scan Code</p>
      <ButtonContainer>
        <IconButton
          text={'Reveal Connection Code'}
          endingIcon={'key'}
          width={250}
          height={48}
          style={{ marginTop: 20 }}
          hovercolor={buttonColor === 'primary' ? '#5881F8' : '#3CBE88'}
          activecolor={buttonColor === 'primary' ? '#5078F2' : '#2FB379'}
          shadowcolor={
            buttonColor === 'primary' ? 'rgba(97, 138, 255, 0.5)' : 'rgba(73, 201, 152, 0.5)'
          }
          onClick={() => getConnectionCode()}
          color={buttonColor}
        />
        <DirectionWrap>
          <IconButton
            text={'Back'}
            width={210}
            height={48}
            buttonType={'text'}
            style={{ color: '#83878b', marginTop: '20px', textDecoration: 'none' }}
            onClick={(e: any) => {
              e.stopPropagation();
              setStep(step - 1);
            }}
            textStyle={{
              fontSize: '15px',
              fontWeight: '500',
              color: '#5F6368'
            }}
            iconStyle={{
              top: '14px'
            }}
            color={buttonColor}
          />
          <IconButton
            text={'Sign in'}
            width={210}
            height={48}
            buttonType={'text'}
            style={{ color: '#83878b', marginTop: '20px', textDecoration: 'none' }}
            onClick={(e: any) => {
              e.stopPropagation();
              setStep(3);
            }}
            textStyle={{
              fontSize: '15px',
              fontWeight: '500',
              color: '#5F6368'
            }}
            iconStyle={{
              top: '14px'
            }}
            color={buttonColor}
          />
        </DirectionWrap>
      </ButtonContainer>
    </>
  );

  const StepThree = () => (
    <ButtonContainer>
      <IconButton
        text={'Sign in'}
        endingIcon={'arrow_forward'}
        width={210}
        height={48}
        style={{ marginTop: 0 }}
        hovercolor={buttonColor === 'primary' ? '#5881F8' : '#3CBE88'}
        activecolor={buttonColor === 'primary' ? '#5078F2' : '#2FB379'}
        shadowcolor={
          buttonColor === 'primary' ? 'rgba(97, 138, 255, 0.5)' : 'rgba(73, 201, 152, 0.5)'
        }
        onClick={(e: any) => {
          e.stopPropagation();
          closeModal();
          ui.setShowSignIn(true);
        }}
        color={buttonColor}
      />

      <IconButton
        text={'Back'}
        width={210}
        height={48}
        buttonType={'text'}
        style={{ color: '#83878b', marginTop: '20px', textDecoration: 'none' }}
        onClick={(e: any) => {
          e.stopPropagation();
          setStep(step - 1);
        }}
        textStyle={{
          fontSize: '15px',
          fontWeight: '500',
          color: '#5F6368'
        }}
        iconStyle={{
          top: '14px'
        }}
        color={buttonColor}
      />
    </ButtonContainer>
  );

  return (
    <>
      <EuiOverlayMask>
        <EuiModal
          onClose={(e: any) => {
            e?.stopPropagation();
            closeModal();
          }}
          style={{
            background: '#F2F3F5',
            padding: '30px',
            borderRadius: '8px',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            minHeight: '585px',
            maxHeight: '585px',
            width: '425px'
          }}
        >
          {step === 1 ? <StepOne /> : step === 2 ? <StepTwo /> : <StepThree />}
        </EuiModal>
      </EuiOverlayMask>
    </>
  );
};

export default observer(StartUpModal);
