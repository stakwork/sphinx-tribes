import { ElementProps } from 'components/interfaces';
import React from 'react';
import styled from 'styled-components';
const T = styled.div`
  color: #3c3f41;
  text-align: left;
  margin-bottom: 10px;

  font-size: 17px;
  font-style: normal;
  font-weight: 500;
  line-height: 23px;
  letter-spacing: 0px;
  color: #3e3f41;
`;
const D = styled.div`
  font-family: Roboto;
  font-size: 10px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 1px;
  text-align: left;
  color: #b0b7bc;
  margin-bottom: 10px;
`;
const P = styled.div`
  font-family: Roboto;
  font-size: 14px;
  color: #5f6368;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0px;
  text-align: left;
  margin-bottom: 10px;
`;

const L = styled.div`
  font-family: Roboto;
  font-size: 13px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0px;
  text-align: left;
  color: #618aff;
`;

function Title(props: ElementProps) {
  return <T {...props}>{props.children}</T>;
}
function Date(props: ElementProps) {
  return <D {...props}>{props.children}</D>;
}
function Paragraph(props: ElementProps) {
  return <P {...props}>{props.children}</P>;
}
function Link(props: ElementProps) {
  return <L {...props}>{props.children}</L>;
}

export { Title, Date, Paragraph, Link };
