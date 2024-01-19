import React from 'react';
import styled from 'styled-components';
import { TextInputProps } from 'components/interfaces';
import { colors } from '../../config/colors';
import { FieldEnv, FieldText } from '../form/inputs/index';

const F = styled((props: any) => <FieldEnv {...props} />)`
  .euiFormLabel[for] {
    cursor: default;
  }
`;

const R = styled.div`
  position: relative;
`;
export default function TextInput({
  label,
  value,
  onChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  style
}: TextInputProps) {
  const color = colors['light'];
  return (
    <>
      <F label={label}>
        <R>
          <FieldText
            color={color}
            name="first"
            value={value || ''}
            readOnly={readOnly || false}
            onChange={(e: any) => onChange(e.target.value)}
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
