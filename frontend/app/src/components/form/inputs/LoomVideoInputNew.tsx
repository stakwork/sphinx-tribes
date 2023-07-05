import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiText } from '@elastic/eui';
import { colors } from '../../../config/colors';
import LoomViewerRecorderNew from '../../../people/utils/LoomViewerRecorderNew';
import type { Props } from './propsType';

interface styleProps {
  color?: any;
  isVideo?: boolean;
}

const LoomVideoContainer = styled.div<styleProps>`
  width: 292px;
  height: 175px;
  left: 698px;
  top: 313px;
  background: ${(p: any) => !p.isVideo && "url('/static/loom_video_outer_border.svg')"};
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
    color: ${(p: any) => p.color && p.color.grayish.G300};
    margin-top: 6px;
  }
`;
export default function LoomVideoInputNew({
  value,
  handleChange,
  handleBlur,
  handleFocus,
  style = {}
}: Props) {
  const color = colors['light'];

  const [isVideo, setIsVideo] = useState<boolean>(false);

  return (
    <LoomVideoContainer color={color} isVideo={isVideo} style={style}>
      <LoomViewerRecorderNew
        name="loomVideo"
        onChange={(e: any) => {
          handleChange(e);
        }}
        loomEmbedUrl={value}
        onBlur={handleBlur}
        onFocus={handleFocus}
        setIsVideo={setIsVideo}
        style={{ marginTop: 59 }}
      />
      {!isVideo && <EuiText className="optionalText">Optional</EuiText>}
    </LoomVideoContainer>
  );
}
