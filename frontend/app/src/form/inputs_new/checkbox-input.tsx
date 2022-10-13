import React, { useEffect } from 'react';
import styled from 'styled-components';
import { EuiCheckbox } from '@elastic/eui';
import type { Props } from './propsType';
import { Note } from './index';

export default function CheckBox({
  label,
  note,
  value,
  name,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  extraHTML
}: Props) {
  useEffect(() => {
    // if value not initiated, default value true
    if (name === 'show' && value === undefined) handleChange(true);
  }, []);

  return (
    <>
      <EuiCheckbox
        id="hi"
        label=""
        checked={value}
        onChange={(e) => {
          handleChange(e.target.checked);
        }}
        onBlur={handleBlur}
        onFocus={handleFocus}
        compressed
      />
      {note && <Note>*{note}</Note>}
    </>
  );
}

const ExtraText = styled.div`
  color: #ddd;
  padding: 10px 10px 25px 10px;
  max-width: calc(100% - 20px);
  word-break: break-all;
  font-size: 14px;
`;
