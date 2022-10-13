import React from 'react';
import styled from 'styled-components';
import * as I from '../interfaces';

export default function Wanted(props: I.Wanted) {
  return (
    <Wrap>
      <div>{props.title}</div>
      <Row>
        <div>{props.priceMin}</div>
        <div> ~ </div>
        <div>{props.priceMax}</div>
      </Row>
    </Wrap>
  );
}

const Wrap = styled.div``;

const Row = styled.div`
  color: #f1f1f1;
  display: flex;
  font-size: 12px;
`;
