import React from 'react';
import styled from 'styled-components';
import type { Props } from './propsType';
import { colors } from '../../colors';
import { EuiText } from '@elastic/eui';
import LoomViewerRecorderNew from '../../people/utils/loomViewerRecorderNew';

export default function LoomVideoInputNew({
  name,
  error,
  note,
  label,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  extraHTML
}: Props) {
  const color = colors['light'];
  return (
    <LoomVideoContainer
      color={color}
      style={{
        marginTop: '55px'
      }}>
      <LoomViewerRecorderNew
        name="loomVideo"
        onChange={(e) => {
          handleChange(e);
        }}
        loomEmbedUrl={value}
        onBlur={handleBlur}
        onFocus={handleFocus}
        style={{ marginTop: 59 }}
      />
      <EuiText className="optionalText">Optional</EuiText>
    </LoomVideoContainer>
  );
}

interface styleProps {
  color?: any;
}

const LoomVideoContainer = styled.div<styleProps>`
  width: 292px;
  height: 175px;
  left: 698px;
  top: 313px;
  background: url('/static/loom_video_outer_border.svg');
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
  .optionalText {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 35px;
    display: flex;
    align-items: center;
    text-align: center;
    color: #b0b7bc;
    margin-top: 6px;
  }
`;
