import React from 'react';
import styled from 'styled-components';
import { FieldEnv, FieldText } from './../form/inputs/index';

export default function TextInput({
  label,
  value,
  onChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  style
}: any) {
  return (
    <>
      <F label={label}>
        <R>
          <FieldText
            name="first"
            value={value || ''}
            readOnly={readOnly || false}
            onChange={(e) => onChange(e.target.value)}
            onBlur={handleBlur}
            onFocus={handleFocus}
            prepend={prepend}
            style={style}
          />
        </R>
      </F>
    </>
  );
}

const F = styled(FieldEnv)`
  .euiFormLabel[for] {
    cursor: default;
  }
`;

const R = styled.div`
  position: relative;
`;
