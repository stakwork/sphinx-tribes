import React from 'react';
import styled from 'styled-components';

const Wrap = styled.div`
  display: flex;
  align-items: center;
`;

interface IconProps {
  source: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p: any) => `url(${p.source})`};
  width: 20px;
  height: 20px;
  margin-right: 5px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: contain; /* Resize the background image to cover the entire container */
  border-radius: 5px;
  overflow: hidden;
`;

export default function TwitterView(props: { handle: string }) {
  return (
    <Wrap>
      <Icon source={`/static/twitter.png`} />
      <div>@{props.handle}</div>
    </Wrap>
  );
}
