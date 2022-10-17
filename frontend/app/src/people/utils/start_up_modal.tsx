import { EuiModal, EuiOverlayMask, EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import IconButton from '../../sphinxUI/icon_button';
import { useStores } from '../../store';

const StartUpModal = ({ closeModal, dataObject, buttonColor }) => {
  const [count, setCount] = useState(0);
  const { ui } = useStores();

  useEffect(() => {
    const interval = setInterval(() => {
      if (count < dataObject.length) {
        setCount((p) => (p + 1) % dataObject.length);
      }
    }, 3000);
    return () => {
      clearInterval(interval);
    };
  }, []);

  return (
    <>
      <EuiOverlayMask>
        <EuiModal
          onClose={closeModal}
          style={{
            background: '#F2F3F5',
            padding: '30px',
            borderRadius: '8px'
          }}>
          <ModalContainer>
            <figure
              style={{
                width: '100%',
                display: 'flex',
                justifyContent: 'center'
              }}>
              {dataObject.map(({ img, step, heading }, index) => {
                return (
                  <Container className={count === index ? 'active' : ''}>
                    <TextContainer>
                      <EuiText
                        style={{
                          fontSize: '14px',
                          fontWeight: '600',
                          color: buttonColor === 'primary' ? '#6189fb' : '#58c998'
                        }}>
                        {step ? step : ''}
                      </EuiText>
                    </TextContainer>
                    <TextContainer>
                      <EuiText
                        style={{
                          fontSize: '28px',
                          fontWeight: '900'
                        }}>
                        {' '}
                        {heading ? heading : ''}{' '}
                      </EuiText>
                    </TextContainer>
                    <img src={img} alt="" height={'164px'} width={'198px'} />
                  </Container>
                );
              })}
            </figure>
          </ModalContainer>

          <ButtonContainer>
            <IconButton
              text={'Get Sphinx'}
              endingIcon={'arrow_forward'}
              width={210}
              height={48}
              style={{ marginTop: 20 }}
              onClick={() => {
                window.open('https://sphinx.chat/', '_blank');
              }}
              color={buttonColor}
            />

            <IconButton
              text={'I have Sphinx'}
              endingIcon={'arrow_forward'}
              width={210}
              height={48}
              buttonType={'text'}
              style={{ color: '#83878b', marginTop: '10px' }}
              onClick={() => {
                closeModal();
                ui.setShowSignIn(true);
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

const Container = styled.div`
  position: absolute;
  opacity: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  &.active {
    opacity: 1;
  }
  transition: all ease 0.3s;
`;

const ModalContainer = styled.div`
  min-height: 300px;
`;

const TextContainer = styled.div`
  width: 100%;
  padding: 10px;
  display: flex;
  justify-content: center;
`;

const ButtonContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;
