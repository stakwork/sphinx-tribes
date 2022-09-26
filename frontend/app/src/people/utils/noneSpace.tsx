import React from 'react';
import styled from 'styled-components';
import { Button } from '../../sphinxUI';

export default function NoneSpace(props) {
  if (props.banner) {
    return (
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100%',
          background: '#f0f1f3',
          padding: 40,
          width: '100%',
          ...props.style
        }}
      >
        <Icon src={`/static/${props.img}`} style={{ width: 180, height: 180 }} />

        <div style={{ marginLeft: 20, padding: 20 }}>
          <H small={props.small} style={{ paddingLeft: 0 }}>
            {props.text}
          </H>
          <C style={{ paddingLeft: 0 }}>{props.sub}</C>

          {props.buttonText && (
            <Button
              text={props.buttonText}
              endingIcon={props.buttonIcon}
              width={210}
              height={48}
              style={{ marginTop: 20 }}
              onClick={props.action}
              color="primary"
            />
          )}
        </div>
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
      <Icon src={`/static/${props.img}`} />

      <>
        <H small={props.small}>{props.text}</H>
        <C>{props.sub}</C>
      </>

      <div style={{ height: 200 }}>
        {props.buttonText && (
          <Button
            text={props.buttonText}
            leadingIcon={props.buttonIcon}
            width={210}
            height={48}
            style={{ marginTop: 40 }}
            onClick={props.action}
            color="primary"
          />
        )}
      </div>
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

interface HProps {
  small: boolean;
}

const H = styled.div<HProps>`
  margin-top: 10px;

  font-family: Roboto;
  font-style: normal;

  display: flex;
  align-items: center;
  text-align: center;

  /* Primary Text 1 */

  color: #292c33;
  padding: 0 10px;
  max-width: 350px;
  letter-spacing: 0px;
  color: rgb(60, 63, 65);

  font-weight: 700;
  font-size: ${(p) => (p.small ? '22px' : '30px')};
  line-height: ${(p) => (p.small ? '26px' : '40px')}; ;
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
