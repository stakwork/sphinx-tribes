import React from 'react';
import type { Props } from './propsType';
import LoomViewerRecorder from '../../people/utils/loomViewerRecorder';
import { colors } from '../../colors';

export default function LoomVideoInput({
  value,
  handleChange,
  handleBlur,
  handleFocus,
}: Props) {
  const color = colors['light'];
  return (
    <>
      <LoomViewerRecorder
        name="loomVideo"
        onChange={(e) => {
          handleChange(e);
        }}
        loomEmbedUrl={value}
        onBlur={handleBlur}
        onFocus={handleFocus}
        style={{ marginBottom: 10 }}
      />
    </>
  );
}
