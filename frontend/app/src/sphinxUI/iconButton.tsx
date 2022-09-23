import React from 'react';
import styled from 'styled-components';
import { EuiButton } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';

export default function IconButton(props: any) {
  const colors = {
    primary: {
      background: '#618AFF',
      color: '#fff'
    },
    white: {
      background: '#fff',
      color: '#5F6368',
      border: '1px solid #DDE1E5'
    },
    link: {
      background: '#fff',
      color: '#618AFF',
      border: '1px solid #A3C1FF'
    }
  };

  return (
    <B
      background={props.style && props.style.background}
      width={props.style && props.style.width}
      style={{ ...colors[props.color], ...props.style }}
      disabled={props.disabled}
      onClick={props.onClick}
    >
      <div style={{ display: 'flex', height: '100%', alignItems: 'center' }}>
        <MaterialIcon
          icon={props.icon}
          style={{ fontSize: props.size ? props.size : 30, color: '#B0B7BC', ...props.iconStyle }}
        />
      </div>
    </B>
  );
}

interface BProps {
  background: string;
  width: string;
}
const B = styled(EuiButton)<BProps>`
  background: ${(p) => (p.background ? p.background : '#ffffff00')} !important;
  position: relative;
  width: fit-content !important;
  min-width: ${(p) => (p.width ? p.width : 'fit-content')};
  max-width: ${(p) => (p.width ? p.width : 'fit-content')};
  width: ${(p) => (p.width ? p.width : '30px')};
  font-weight: bold;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: none !important;
  text-transform: none !important;

  .euiButton__content {
    height: 100%;
    width: 100%;
    vertical-align: middle;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0;
  }
`;
