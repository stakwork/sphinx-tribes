import React from 'react';
import styled from 'styled-components';
import * as I from '../interfaces';

const Wrap = styled.div`
  display: flex;
  width: 100%;
`;

const Sub = styled.div`
  color: #f1f1f1;
  display: flex;
  font-size: 12px;
`;
export default function Offer(props: I.Offer) {
  return (
    <Wrap>
      <div>
        <div>{props.title}</div>
        <Sub>
          <div>{props.price}</div>
        </Sub>
      </div>
    </Wrap>
  );
}
