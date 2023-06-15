import React from 'react';
import LoomViewerRecorder from '../../../people/utils/loomViewerRecorder';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';

export default function LoomVideoInput({ value, handleChange, handleBlur, handleFocus }: Props) {
  const color = colors['light'];
  return (
    <>
      <LoomViewerRecorder
        name="loomVideo"
        onChange={(e: any) => {
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
