import React from 'react';
import styled from 'styled-components';
import { EuiLoadingSpinner } from '@elastic/eui';

export default function PageLoadSpinner(props) {
  if (props.noAnimate) {
    return (
      <BottomLoadmoreWrap show={props.show} style={props.style}>
        <EuiLoadingSpinner size="l" style={{ marginLeft: -10, padding: 10 }} />
      </BottomLoadmoreWrap>
    );
  }

  return (
    <LoadmoreWrap show={props.show} style={props.style}>
      <EuiLoadingSpinner size="l" style={{ marginLeft: -10, padding: 10 }} />
    </LoadmoreWrap>
  );
}

interface LoadmoreWrapProps {
  show: boolean;
}
export const LoadmoreWrap = styled.div<LoadmoreWrapProps>`
  position: relative;
  text-align: center;
  width: 100%;
  // background:#ffffff33;
  min-width: 100%;
  transition: all 0.2s;
  display: flex;
  justify-content: center;
  height: ${(p) => (p.show ? '50px' : '0px')};
  opacity: ${(p) => (p.show ? '1' : '0')};
  overflow: hidden;
  visibility: ${(p) => (p.show ? 'visible' : 'hidden')};
`;
export const BottomLoadmoreWrap = styled.div<LoadmoreWrapProps>`
  position: relative;
  text-align: center;
  width: 100%;
  min-width: 100%;
  height: 50px;
  visibility: ${(p) => (p.show ? 'visible' : 'hidden')};
`;
