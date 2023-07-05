import React from 'react';
import styled from 'styled-components';
import { NoneSpaceProps } from 'people/interfaces';
import IconButton from '../../components/common/IconButton2';

const H = styled.div`
  margin-top: 10px;

  font-family: Roboto;
  font-style: normal;

  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;

  /* Primary Text 1 */

  color: #292c33;
  padding: 0 10px;
  letter-spacing: 0px;
  color: rgb(60, 63, 65);

  font-weight: 700;
`;

const C = styled.div`
  margin-top: 10px;
  font-family: Roboto;
  font-size: 22px;
  font-style: normal;
  font-weight: 400;
  line-height: 26px;
  letter-spacing: 0em;
  text-align: center;
  color: #8e969c;

  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 18px;
  text-align: center;

  /* Main bottom icons */

  color: #5f6368;
  padding: 0 10px;

  max-width: 350px;
  padding: 0 65px;
`;

const ButtonContainer = styled.div`
  display: flex;
  justify-content: center;
  flex-direction: row;
`;

const ButtonContainerMobile = styled.div`
  display: flex;
  flex-direction: column;
`;
export default function NoneSpaceHomePage(props: NoneSpaceProps) {
  if (props.banner) {
    return (
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100%',
          background: '#f0f1f3',
          padding: 40,
          width: '100%',
          ...props.style
        }}
      >
        <H style={{ paddingLeft: 0, fontSize: '100px', fontFamily: 'Barlow' }}>{props.text}</H>
        <C style={{ paddingLeft: 0, fontFamily: 'Barlow' }}>{props.sub}</C>

        <ButtonContainer>
          {props.buttonText1 && (
            <IconButton
              text={props.buttonText1}
              endingIcon={props.buttonIcon}
              width={210}
              height={48}
              style={{ marginTop: 20 }}
              onClick={props.action1}
              color="primary"
              hovercolor={'#5881F8'}
              activecolor={'#5078F2'}
              shadowcolor={'rgba(97, 138, 255, 0.5)'}
              iconStyle={{
                top: '13px',
                right: '14px'
              }}
            />
          )}

          {props.buttonText2 && (
            <IconButton
              text={props.buttonText2}
              endingIcon={props.buttonIcon}
              width={210}
              height={48}
              style={{ marginTop: 20, marginLeft: 10 }}
              onClick={props.action2}
              color="success"
              hovercolor={'#3CBE88'}
              activecolor={'#2FB379'}
              shadowcolor={'rgba(73, 201, 152, 0.5)'}
            />
          )}
        </ButtonContainer>
      </div>
    );
  }

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100%',
        background: '#f0f1f3',
        ...props.style
      }}
    >
      <>
        <H
          style={{
            fontSize: '40px',
            fontFamily: 'Barlow'
          }}
        >
          {props.text}
        </H>
        <C>{props.sub}</C>
      </>
      <ButtonContainerMobile>
        {props.buttonText1 && (
          <IconButton
            text={props.buttonText1}
            endingIcon={props.buttonIcon}
            width={210}
            height={48}
            style={{ marginTop: 40 }}
            onClick={props.action1}
            color="primary"
            hovercolor={'#5881F8'}
            activecolor={'#5078F2'}
            shadowcolor={'rgba(97, 138, 255, 0.5)'}
          />
        )}
        {props.buttonText2 && (
          <IconButton
            text={props.buttonText2}
            endingIcon={props.buttonIcon}
            width={210}
            height={48}
            style={{ marginTop: 20 }}
            onClick={props.action2}
            color="success"
            hovercolor={'#5881F8'}
            activecolor={'#5078F2'}
            shadowcolor={'rgba(97, 138, 255, 0.5)'}
          />
        )}
      </ButtonContainerMobile>
    </div>
  );
}
