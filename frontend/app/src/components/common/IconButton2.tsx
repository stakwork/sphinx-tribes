import React from 'react';
import styled from 'styled-components';
import { EuiButton, EuiLoadingSpinner } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { IconButtonProps } from 'components/interfaces';

interface ButtonHoverProps {
  hovercolor?: string;
  activecolor?: string;
  shadowcolor?: string;
}

const B = styled(EuiButton)<ButtonHoverProps>`
  position: relative;
  border-radius: 100px;
  height: 36px;
  width: 36px;
  min-width: 36px;
  font-weight: bold;
  border: none;
  font-weight: 500;
  font-size: 15px;
  font-family: 'Barlow';
  line-height: 18px;
  display: flex;
  align-items: center;
  text-align: center;
  box-shadow: none !important;
  text-transform: none !important;
  transform: none !important;
  text-decoration: none !important;
  box-shadow: ${(p: any) => (p.shadowcolor ? `0px 2px 10px ${p.shadowcolor}` : 'none')} !important;

  &:hover {
    background: ${(p: any) => (p.hovercolor ? p.hovercolor : undefined)} !important;
    transform: none !important;
    text-decoration: none !important;
  }

  &:active {
    background: ${(p: any) => (p.activecolor ? p.activecolor : 'none')} !important;
    transform: none !important;
    text-decoration: none !important;
  }
  @media only screen and (max-width: 700px) {
    font-size: 12.5px;
  }
`;

const T = styled(EuiButton)`
  position: relative;
  border-radius: 100px;
  height: 36px;
  border: none;

  font-weight: 500;
  font-size: 15px;
  font-family: 'Barlow';
  line-height: 18px;
  display: flex;
  align-items: center;
  text-align: center;
  box-shadow: none !important;
  text-transform: none !important;
  transform: none !important;
  text-decoration: none !important;

  &.textButton {
    background: transparent;
    transform: none !important;
    text-decoration: none !important;
  }
  &.textButton:hover {
    background: transparent;
    transform: none !important;
    text-decoration: none !important;
  }
  &.textButton:focus {
    background: transparent !important;
    transform: none !important;
    text-decoration: none !important;
  }
`;
function hexToRgba(hex: string, opacity: any = 1) {
  try {
    const shorthandRegex = /^#?([a-f\d])([a-f\d])([a-f\d])$/i;
    hex = hex.replace(shorthandRegex, function (m: any, r: any, g: any, b: any) {
      return r + r + g + g + b + b;
    });

    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    const rgb = result
      ? {
          r: parseInt(result[1], 16),
          g: parseInt(result[2], 16),
          b: parseInt(result[3], 16)
        }
      : null;

    return rgb ? `rgba(${rgb?.r}, ${rgb?.g}, ${rgb?.b}, ${opacity})` : undefined;
  } catch {
    return undefined;
  }
}

export default function IconButton(props: IconButtonProps) {
  const { iconStyle, id, color = 'primary' } = props;

  const colors = {
    primary: {
      background: '#618AFF',
      color: '#fff'
    },
    success: {
      background: '#49C998',
      color: '#fff'
    },
    error: {
      background: '#ED7474',
      color: '#fff'
    },
    white: {
      background: '#fff',
      color: '#5F6368',
      border: '1px solid #DDE1E5'
    },
    clear: {
      background: '#fff',
      color: '#5F6368',
      border: '1px solid #fff'
    },
    noColor: {
      background: '#ffffff00',
      color: '#5F6368'
    },
    link: {
      background: '#fff',
      color: '#618AFF',
      border: '1px solid #A3C1FF'
    },
    link2: {
      background: '#fff',
      color: '#618AFF',
      width: 'fit-content',
      minWidth: 'fit-content',
      border: 'none',
      padding: 0,
      margin: 0
    },
    widget: {
      background: '#DDE1E5',
      color: '#3C3F41'
    },
    danger: {
      background: 'red',
      color: '#ffffff'
    },
    desktopWidget: {
      background: 'rgba(0,0,0,0)',
      color: '#3C3F41',
      border: '1px dashed #B0B7BC',
      boxSizing: 'border-box',
      borderRadius: 4
    },
    transparent: {
      background: '#ffffff12',
      color: '#fff',
      border: 'none'
    }
  } as const;

  return props.buttonType && props.buttonType === 'text' ? (
    <T
      className="textButton"
      id={id}
      style={{
        padding: props.icon && '0 0 0 15px',
        position: 'relative',
        opacity: props.disabled ? 0.7 : 1,
        height: props.height,
        width: props.width,
        paddingRight: props.leadingIcon && 10,
        ...props.style
      }}
      disabled={props.disabled}
      onClick={props.onClick}
    >
      <span style={{ ...props.textStyle }}>
        {props.loading || props.submitting ? <EuiLoadingSpinner size="m" /> : <>{props.text}</>}
      </span>
      {props.endingIcon && (
        <MaterialIcon
          icon={props.endingIcon}
          style={{
            fontSize: props.iconSize ? props.iconSize : 20,
            marginLeft: 10,
            position: 'absolute',
            right: '20px',
            ...iconStyle
          }}
        />
      )}
    </T>
  ) : (
    <B
      id={id}
      style={{
        ...colors[color],
        padding: props.icon && '0 0 0 15px',
        position: 'relative',
        opacity: props.disabled ? 0.7 : 1,
        height: props.height,
        width: props.width,
        paddingRight: props.leadingIcon && 10,
        ...props.style
      }}
      hovercolor={props.hovercolor ?? hexToRgba(colors[color].background, 0.8)}
      activecolor={props.activecolor}
      shadowcolor={props.shadowcolor}
      disabled={props.disabled}
      className="test"
      onClick={props.onClick}
    >
      <div>
        {props.icon && (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              position: 'absolute',
              top: 0,
              left: 3,
              height: '100%',
              ...iconStyle
            }}
          >
            <MaterialIcon
              icon={props.icon}
              style={{ fontSize: props.iconSize ? props.iconSize : 30, ...iconStyle }}
            />
          </div>
        )}

        <div style={{ display: 'flex', alignItems: 'center' }}>
          {props.leadingIcon && (
            <MaterialIcon
              icon={props.leadingIcon}
              style={{
                fontSize: props.iconSize ? props.iconSize : 20,
                marginRight: 10,
                ...iconStyle
              }}
            />
          )}
          {props.leadingImg && (
            <div
              style={{
                ...props.leadingImgStyle
              }}
            >
              <img height={'100%'} width={'100%'} src={props.leadingImg} alt="leading" />
            </div>
          )}
          <span style={{ ...props.textStyle }}>
            {props.loading || props.submitting ? <EuiLoadingSpinner size="m" /> : <>{props.text}</>}
          </span>
          {props.endingIcon && (
            <MaterialIcon
              icon={props.endingIcon}
              style={{
                fontSize: props.iconSize ? props.iconSize : 20,
                marginLeft: 10,
                position: 'absolute',
                right: '20px',
                ...iconStyle
              }}
            />
          )}
          {props.endingImg && (
            <div
              style={{
                ...props.endingImgStyle
              }}
            >
              <img height={'100%'} width={'100%'} src={props.endingImg} alt="leading" />
            </div>
          )}
        </div>
      </div>
    </B>
  );
}
