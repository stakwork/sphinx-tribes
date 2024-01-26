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
  width: 4.6875rem;
  height: 4.728rem;
`;

const Icon = styled.div`
  width: 4.6875rem;
  height: 4.728rem;
  display: flex;
  align-items: center;
`;
export default function QR(props: QRProps) {
  const { type } = props;
  const color = colors['light'];

  const centerIcon =
    type === 'connect' ? (
      <Icon data-testid={'testid-connectimg'}>
        <MaterialIcon icon={'person_add'} style={{ fontSize: 36, marginLeft: 7 }} />
      </Icon>
    ) : (
      <Img src={'/static/sphinx.svg'} />
    );

  return (
    <div style={{ position: 'relative' }}>
      <QRCode
        data-testid={'testid-qrcode'}
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

      {type === 'connect' && (
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
      )}
    </div>
  );
}
