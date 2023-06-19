import React from 'react';
import styled from 'styled-components';
import { QRCode } from 'react-qr-svg';
import MaterialIcon from '@material/react-material-icon';
import { QRProps } from 'people/interfaces';
import { colors } from '../../config/colors';

interface ImageProps {
  readonly src: string;
}
const Img = styled.div<ImageProps>`
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  height: 55px;
  width: 55px;
  border-radius: 50%;
`;

const Icon = styled.div`
  height: 55px;
  width: 55px;
  border-radius: 50%;
  display: flex;
  align-items: center;
`;
export default function QR(props: QRProps) {
  const { type } = props;
  const color = colors['light'];

  const centerIcon =
    type === 'connect' ? (
      <Icon>
        <MaterialIcon icon={'person_add'} style={{ fontSize: 36, marginLeft: 7 }} />
      </Icon>
    ) : (
      <Img src={'/static/sphinx.png'} />
    );

  return (
    <div style={{ position: 'relative' }}>
      <QRCode
        bgColor={color.pureWhite}
        fgColor={color.pureBlack}
        level={'Q'}
        style={{ width: props.size }}
        value={props.value}
      />

      {/* logo env */}
      <div
        style={{
          position: 'absolute',
          zIndex: 10,
          height: props.size,
          width: props.size,
          top: 0,
          left: 0,
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center'
        }}
      >
        {centerIcon}
      </div>

      <div
        style={{
          position: 'absolute',
          zIndex: 8,
          height: props.size,
          width: props.size,
          top: 0,
          left: 0,
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center'
        }}
      >
        <div style={{ background: color.pureWhite, height: 75, width: 75 }} />
      </div>
    </div>
  );
}
