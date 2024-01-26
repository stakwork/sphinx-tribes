import React from 'react';
import styled from 'styled-components';
import { EuiButton, EuiLoadingSpinner } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import { ButtonProps } from 'components/interfaces';
import { commonColors } from '../../config/commonColors';

const B = styled(EuiButton)`
  position: relative;
  border-radius: 100px;
  height: 36px;
  font-weight: bold;
  border: none;
  font-weight: 500;
  font-size: 14px;
  line-height: 18px;
  display: flex;
  align-items: center;
  text-align: center;
  box-shadow: none !important;
  text-transform: none !important;
  text-decoration: none !important;
  transform: none !important;
  @media only screen and (max-width: 700px) {
    font-size: 12.5px;
  }
  @media only screen and (max-width: 500px) {
    font-size: 11.5px;
  }
`;

interface IconProps {
  src: string;
}

const Img = styled.div<IconProps>`
  background-image: ${(p: any) => `url(${p.src})`};
  width: 80px;
  height: 80px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: cover; /* Resize the background image to cover the entire container */
  border-radius: 80px;
  overflow: hidden;
`;
export default function Button(props: ButtonProps) {
  const { iconStyle, id } = props;

  return (
    <B
      id={id}
      style={{
        ...commonColors[props.color ? props.color : 'primary'],
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
      {props.children ? (
        props.children
      ) : props.wideButton ? (
        <div
          style={{
            display: 'flex',
            justifyContent: 'flex-start',
            width: '100%',
            minWidth: '100%'
          }}
        >
          {props.text && (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                position: 'absolute',
                left: 18,
                top: 0,
                height: props.height,
                maxWidth: '80%'
              }}
            >
              <div
                style={{
                  overflow: 'hidden',
                  whiteSpace: 'nowrap',
                  textOverflow: 'ellipsis'
                }}
              >
                {props.text}
              </div>
            </div>
          )}
          {props.icon && (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                position: 'absolute',
                right: 13,
                top: 0,
                height: props.height
              }}
            >
              <MaterialIcon
                icon={props.icon}
                style={{ fontSize: props.iconSize ? props.iconSize : 30, ...iconStyle }}
              />
            </div>
          )}
        </div>
      ) : (
        <>
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
          {props.img && (
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                position: 'absolute',
                top: 0,
                left: 10,
                height: '100%'
              }}
            >
              <Img
                src={`/static/${props.img}`}
                style={{
                  width: props.imgSize ? props.imgSize : 30,
                  height: props.imgSize ? props.imgSize : 30,
                  ...props.imgStyle
                }}
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

            {props.leadingImgUrl && (
              <Img
                src={props.leadingImgUrl}
                style={{
                  width: props.imgSize ? props.imgSize : 28,
                  height: props.imgSize ? props.imgSize : 28,
                  ...props.imgStyle
                }}
              />
            )}

            {props.loading || props.submitting ? (
              <EuiLoadingSpinner size="m" />
            ) : (
              <div
                style={{
                  ...props.ButtonTextStyle
                }}
              >
                {props.text}
              </div>
            )}
            {props.endingIcon && (
              <MaterialIcon
                icon={props.endingIcon}
                style={{
                  fontSize: props.iconSize ? props.iconSize : 20,
                  marginLeft: 10,
                  ...iconStyle
                }}
              />
            )}
          </div>
        </>
      )}
    </B>
  );
}
