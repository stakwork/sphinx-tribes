import React from 'react';
import styled from 'styled-components';
import * as I from '../interfaces';

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

const Wrap = styled.div`
  display: flex;
  width: 100%;
`;

const Sub = styled.div`
  color: #f1f1f1;
  display: flex;
  font-size: 12px;
`;

const Image = styled.div`
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  height: 42px;
  width: 70px;
  margin-right: 5px;
`;
