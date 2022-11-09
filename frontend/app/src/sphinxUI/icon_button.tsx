import React from 'react';
import styled from 'styled-components';
import { EuiButton, EuiLoadingSpinner } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';

export default function IconButton(props: any) {
  const { iconStyle, id } = props;

  const colors = {
    primary: {
      background: '#618AFF',
      color: '#fff'
    },
    success: {
      background: '#49C998',
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
  };

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
        ...colors[props.color],
        padding: props.icon && '0 0 0 15px',
        position: 'relative',
        opacity: props.disabled ? 0.7 : 1,
        height: props.height,
        width: props.width,
        paddingRight: props.leadingIcon && 10,
        ...props.style
      }}
      hoverColor={props.hoverColor}
      activeColor={props.activeColor}
      shadowColor={props.shadowColor}
      disabled={props.disabled}
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
        </div>
      </div>
    </B>
  );
}

interface ButtonHoverProps {
  hoverColor?: string;
  activeColor?: string;
  shadowColor?: string;
}

const B = styled(EuiButton)<ButtonHoverProps>`
  position: relative;
  border-radius: 100px;
  height: 36px;
  font-weight: bold;
  border: none;
  font-weight: 500;
  font-size: 15px;
  font-family: Barlow;
  line-height: 18px;
  display: flex;
  align-items: center;
  text-align: center;
  box-shadow: none !important;
  text-transform: none !important;
  transform: none !important;
  text-decoration: none !important;
  box-shadow: ${(p) => (p.shadowColor ? `0px 2px 10px ${p.shadowColor}` : 'none')} !important;

  &:hover {
    background: ${(p) => (p.hoverColor ? p.hoverColor : 'none')} !important;
    transform: none !important;
    text-decoration: none !important;
  }

  &:active {
    background: ${(p) => (p.activeColor ? p.activeColor : 'none')} !important;
    transform: none !important;
    text-decoration: none !important;
  }
`;

const T = styled(EuiButton)`
  position: relative;
  border-radius: 100px;
  height: 36px;
  border: none;

  font-weight: 500;
  font-size: 15px;
  font-family: Barlow;
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
