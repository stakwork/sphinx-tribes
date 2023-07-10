import React from 'react';
import styled from 'styled-components';
import { EuiButton } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { IconButtonProps } from 'components/interfaces';
import { colors } from '../../config/colors';
interface BProps {
  background: string;
  width: string;
}
const B = styled(EuiButton)<BProps>`
  background: ${(p: any) => (p.background ? p.background : '#ffffff00')} !important;
  position: relative;
  width: fit-content !important;
  min-width: ${(p: any) => (p.width ? p.width : 'fit-content')};
  max-width: ${(p: any) => (p.width ? p.width : 'fit-content')};
  width: ${(p: any) => (p.width ? p.width : '30px')};
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

export default function IconButton(props: IconButtonProps) {
  const color = colors['light'];
  const ButtonColors = {
    primary: {
      background: color.blue1,
      color: color.pureWhite
    },
    white: {
      background: color.pureWhite,
      color: color.grayish.G50,
      border: `1px solid ${color.grayish.G600}`
    },
    link: {
      background: color.pureWhite,
      color: color.blue1,
      border: `1px solid ${color.textBlue1}`
    }
  };

  return (
    <B
      background={props.style && props.style.background}
      width={props.style && props.style.width}
      style={{ ...ButtonColors[props.color || 'primary'], ...props.style }}
      disabled={props.disabled}
      onClick={props.onClick}
    >
      <div style={{ display: 'flex', height: '100%', alignItems: 'center' }}>
        <MaterialIcon
          icon={props.icon}
          style={{
            fontSize: props.size ? props.size : 30,
            color: color.grayish.G300,
            ...props.iconStyle
          }}
        />
      </div>
    </B>
  );
}
