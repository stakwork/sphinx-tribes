import React, { useEffect } from 'react';
import { EuiCheckbox } from '@elastic/eui';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { Note } from './index';

export default function CheckBox({
  note,
  value,
  name,
  handleChange,
  handleBlur,
  handleFocus
}: Props) {
  const color = colors['light'];
  useEffect(() => {
    // if value not initiated, default value true
    if (name === 'show' && value === undefined) handleChange(true);
  }, [handleChange, name, value]);

  return (
    <>
      <EuiCheckbox
        id="hi"
        label=""
        checked={value}
        onChange={(e: any) => {
          handleChange(e.target.checked);
        }}
        onBlur={handleBlur}
        onFocus={handleFocus}
        compressed
      />
      {note && <Note color={color}>*{note}</Note>}
    </>
  );
}
