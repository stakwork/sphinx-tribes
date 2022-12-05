import React from 'react';
import styled from 'styled-components';
import type { Props } from './propsType';
import { FieldEnv } from './index';
import LoomViewerRecorder from '../../people/utils/loomViewerRecorder';
import { colors } from '../../colors';

export default function LoomVideoInput({
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
