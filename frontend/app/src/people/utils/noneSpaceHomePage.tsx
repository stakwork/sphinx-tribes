import React from 'react';
import styled from 'styled-components';
import { Button } from '../../sphinxUI';

export default function NoneSpaceHomePage(props) {
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
        }}>
        <H style={{ paddingLeft: 0, fontSize: '100px' }}>{props.text}</H>
        <C style={{ paddingLeft: 0 }}>{props.sub}</C>

        <ButtonContainer>
          {props.buttonText1 && (
            <Button
              text={props.buttonText1}
              endingIcon={props.buttonIcon}
              width={210}
              height={48}
              style={{ marginTop: 20 }}
              onClick={props.action1}
              color="primary"
            />
          )}

          {props.buttonText2 && (
            <Button
              text={props.buttonText2}
              endingIcon={props.buttonIcon}
              width={210}
              height={48}
              style={{ marginTop: 20, marginLeft: 10 }}
              onClick={props.action2}
              color="success"
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
      }}>
      <>
        <H
          style={{
            fontSize: '40px'
          }}>
          {props.text}
        </H>
        <C>{props.sub}</C>
      </>
      <ButtonContainerMobile>
        {props.buttonText1 && (
          <Button
            text={props.buttonText1}
            endingIcon={props.buttonIcon}
            width={210}
            height={48}
            style={{ marginTop: 40 }}
            onClick={props.action1}
            color="primary"
          />
        )}
        {props.buttonText2 && (
          <Button
            text={props.buttonText2}
            endingIcon={props.buttonIcon}
            width={210}
            height={48}
            style={{ marginTop: 20 }}
            onClick={props.action2}
            color="success"
          />
        )}
      </ButtonContainerMobile>
    </div>
  );
}

interface IconProps {
  src: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p) => `url(${p.src})`};
  width: 160px;
  height: 160px;
  margin-right: 10px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: contain; /* Resize the background image to cover the entire container */
  border-radius: 5px;
  overflow: hidden;
`;

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
