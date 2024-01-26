import React from 'react';
import styled from 'styled-components';
import { EuiLoadingSpinner } from '@elastic/eui';
import { PageLoadProps } from 'people/interfaces';

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
  height: ${(p: any) => (p.show ? '50px' : '0px')};
  opacity: ${(p: any) => (p.show ? '1' : '0')};
  overflow: hidden;
  visibility: ${(p: any) => (p.show ? 'visible' : 'hidden')};
`;
export const BottomLoadmoreWrap = styled.div<LoadmoreWrapProps>`
  position: relative;
  text-align: center;
  width: 100%;
  min-width: 100%;
  height: 50px;
  visibility: ${(p: any) => (p.show ? 'visible' : 'hidden')};
`;

export default function PageLoadSpinner(props: PageLoadProps) {
  if (props.noAnimate) {
    return (
      <BottomLoadmoreWrap show={props.show} style={props.style}>
        <EuiLoadingSpinner size="l" style={{ marginLeft: 0, padding: 10 }} />
      </BottomLoadmoreWrap>
    );
  }

  return (
    <LoadmoreWrap show={props.show} style={props.style}>
      <EuiLoadingSpinner size="l" style={{ marginLeft: 0, padding: 10 }} />
    </LoadmoreWrap>
  );
}
