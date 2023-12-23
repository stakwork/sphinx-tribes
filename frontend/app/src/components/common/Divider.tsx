import React from 'react';
import styled from 'styled-components';

const D = styled.div`
  height: 1px;
  background: #ebedef; //#DDE1E5;//#EBEDEF;
  width: 100%;
`;
export default function Divider(props: { style?: React.CSSProperties }) {
  return <D style={{ ...props.style }} data-testid={'testid-divider'} />;
}
