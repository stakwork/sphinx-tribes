/* eslint-disable @typescript-eslint/typedef */
import { Modal, ModalProps } from '@mui/base';
import { Box, styled } from '@mui/system';
import clsx from 'clsx';
import React from 'react';

export type BaseModalProps = ModalProps;

const StyledModal = styled(Modal)`
  position: fixed;
  z-index: 1000000;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
`;

// eslint-disable-next-line react/display-name
const Backdrop = React.forwardRef<HTMLDivElement, { open?: boolean; className: string }>(
  (props, ref) => {
    const { open, className, ...other } = props;
    return <div className={clsx({ 'MuiBackdrop-open': open }, className)} ref={ref} {...other} />;
  }
);

const StyledBackdrop = styled(Backdrop)`
  z-index: -1;
  position: fixed;
  inset: 0;
  background-color: rgb(255 255 255 / 0.8);
  -webkit-tap-highlight-color: transparent;
`;

const Inner = styled('div')(() => ({
  backgroundColor: 'white',
  borderRadius: '0.5rem',
  boxShadow: '0px 4px 20px 0px rgba(0, 0, 0, 0.15)',
  '&:focus, &:focus-visible': {
    outline: 'none'
  },
  '*': {
    fontFamily: 'Barlow'
  }
}));

export const BaseModal = ({ children, ...props }: BaseModalProps) => (
  <StyledModal {...props} slots={{ backdrop: StyledBackdrop }}>
    <Inner>{children}</Inner>
  </StyledModal>
);
