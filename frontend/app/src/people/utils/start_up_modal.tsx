import { EuiModal, EuiOverlayMask } from '@elastic/eui';
import React from 'react';
import styled from 'styled-components';
import IconButton from '../../components/common/icon_button';
import { useStores } from '../../store';
import QR from './QR';
import api from '../../api';

const StartUpModal = ({ closeModal, dataObject, buttonColor }) => {
  const { ui, main } = useStores();

  async function getConnectionCode() {
    if (!ui.connection_string) {
      const code = await api.get('connectioncodes');
      if (code.connection_string) {
        ui.setConnectionString(code.connection_string);
        main.getPeople({ resetPage: true });
      }
    }
  }

  return (
    <>
      <EuiOverlayMask>
        <EuiModal
          onClose={(e) => {
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
          <ModalContainer>
            {ui.connection_string ? (
              <QrContainer>
                <QR size={200} value={ui.connection_string} />
                <QRText>
                  Scan the QR code with Sphinx app to Sign up, then click the Sign in button below
                </QRText>
              </QrContainer>
            ) : (
              <img
                src={
                  dataObject === 'getWork'
                    ? '/static/create_profile_blue.gif'
                    : '/static/create_profile_green.gif'
                }
                height={'274px'}
                alt=""
              />
            )}
          </ModalContainer>
          <ButtonContainer>
            <IconButton
              text={'Get Sphinx'}
              endingIcon={'arrow_forward'}
              width={210}
              height={48}
              style={{ marginTop: 20 }}
              hoverColor={buttonColor === 'primary' ? '#5881F8' : '#3CBE88'}
              activeColor={buttonColor === 'primary' ? '#5078F2' : '#2FB379'}
              shadowColor={
                buttonColor === 'primary' ? 'rgba(97, 138, 255, 0.5)' : 'rgba(73, 201, 152, 0.5)'
              }
              onClick={(e) => {
                e.stopPropagation();
                window.open('https://sphinx.chat/', '_blank');
              }}
              color={buttonColor}
            />
            {ui.connection_string ? (
              <IconButton
                text={'Sign in'}
                endingIcon={'arrow_forward'}
                width={210}
                height={48}
                style={{ marginTop: 20 }}
                hoverColor={buttonColor === 'primary' ? '#5881F8' : '#3CBE88'}
                activeColor={buttonColor === 'primary' ? '#5078F2' : '#2FB379'}
                shadowColor={
                  buttonColor === 'primary' ? 'rgba(97, 138, 255, 0.5)' : 'rgba(73, 201, 152, 0.5)'
                }
                onClick={(e) => {
                  e.stopPropagation();
                  closeModal();
                  ui.setShowSignIn(true);
                }}
                color={buttonColor}
              />
            ) : (
              <IconButton
                text={'Connection Code'}
                endingIcon={'key'}
                width={210}
                height={48}
                style={{ marginTop: 20 }}
                hoverColor={buttonColor === 'primary' ? '#5881F8' : '#3CBE88'}
                activeColor={buttonColor === 'primary' ? '#5078F2' : '#2FB379'}
                shadowColor={
                  buttonColor === 'primary' ? 'rgba(97, 138, 255, 0.5)' : 'rgba(73, 201, 152, 0.5)'
                }
                onClick={getConnectionCode}
                color={buttonColor}
              />
            )}
            <IconButton
              text={'I have Sphinx'}
              endingIcon={'arrow_forward'}
              width={210}
              height={48}
              buttonType={'text'}
              style={{ color: '#83878b', marginTop: '20px', textDecoration: 'none' }}
              onClick={(e) => {
                e.stopPropagation();
                closeModal();
                ui.setShowSignIn(true);
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
        </EuiModal>
      </EuiOverlayMask>
    </>
  );
};

export default StartUpModal;

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
