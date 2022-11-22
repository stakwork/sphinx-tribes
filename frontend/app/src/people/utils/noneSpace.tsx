import React from 'react';
import styled from 'styled-components';
import { colors } from '../../colors';
import { Button } from '../../sphinxUI';

export default function NoneSpace(props) {
  const color = colors['light'];
  if (props.banner) {
    return (
      <OuterContainer
        style={{
          padding: 40,
          width: '100%',
          ...props.style
        }}
        color={color}>
        <Icon src={`/static/${props.img}`} style={{ width: 180, height: 180 }} color={color} />

        <div style={{ marginLeft: 20, padding: 20 }}>
          <H small={props.small} style={{ paddingLeft: 0 }} color={color}>
            {props.text}
          </H>
          <C style={{ paddingLeft: 0 }} color={color}>
            {props.sub}
          </C>

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
      </OuterContainer>
    );
  }

  return (
    <OuterContainer
      style={{
        ...props.style
      }}
      color={color}>
      <Icon src={`/static/${props.img}`} color={color} />

      <>
        <H small={props.small} color={color}>
          {props.text}
        </H>
        <C color={color}>{props.sub}</C>
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
    </OuterContainer>
  );
}

interface IconProps {
  src: string;
  color?: any;
}

interface styledProps {
  color?: any;
}

const OuterContainer = styled.div<styledProps>`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  background: ${(p) => p?.color && p.color.background100};
`;

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
  color?: any;
}

const H = styled.div<HProps>`
  margin-top: 10px;

  font-family: Roboto;
  font-style: normal;
  display: flex;
  align-items: center;
  text-align: center;
  color: ${(p) => p?.color && p.color.grayish.G05};
  padding: 0 10px;
  max-width: 350px;
  letter-spacing: 0px;
  color: ${(p) => p?.color && p.color.grayish.G07};
  font-weight: 700;
  font-size: ${(p) => (p.small ? '22px' : '30px')};
  line-height: ${(p) => (p.small ? '26px' : '40px')}; ;
`;

const C = styled.div<styledProps>`
  margin-top: 10px;
  font-family: Roboto;
  font-size: 22px;
  font-style: normal;
  font-weight: 400;
  line-height: 26px;
  letter-spacing: 0em;
  text-align: center;
  color: ${(p) => p?.color && p.color.grayish.G100};
  font-family: Roboto;
  font-style: normal;
  font-weight: normal;
  font-size: 15px;
  line-height: 18px;
  text-align: center;
  color: ${(p) => p?.color && p.color.grayish.G50};
  padding: 0 10px;
  max-width: 350px;
  padding: 0 65px;
`;
