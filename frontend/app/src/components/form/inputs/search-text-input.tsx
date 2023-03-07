import React from 'react';
import styled from 'styled-components';
import { EuiIcon } from '@elastic/eui';
import { colors } from '../../../config/colors';

export default function SearchTextInput({
  placeholder,
  error,
  label,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
}: any) {
  let labeltext = label;
  if (error) labeltext = `${labeltext} (${error})`;
  const color = colors['light'];

  return (
    <>
      <R>
        <Text
          color={color}
          name="first"
          value={value || ''}
          readOnly={readOnly || false}
          onChange={(e) => handleChange(e.target.value)}
          onBlur={handleBlur}
          onFocus={handleFocus}
          placeholder={placeholder || 'Search'}
        />
        {error && (
          <E color={color}>
            <EuiIcon type="alert" size="m" style={{ width: 20, height: 20 }} />
          </E>
        )}
      </R>
    </>
  );
}

interface styledProps {
  color?: any;
}

const Text = styled.input<styledProps>`
  background: ${(p) => p.color && p.color.grayish.G71A};
  border: 1px solid ${(p) => p.color && p.color.grayish.G750};
  box-sizing: border-box;
  border-radius: 21px;
  padding-left: 20px;
  width: 100%;
  font-style: normal;
  font-weight: normal;
  font-size: 12px;
  line-height: 14px;
  height: 35px;
`;

const E = styled.div<styledProps>`
  position: absolute;
  right: 10px;
  top: 0px;
  display: flex;
  height: 100%;
  justify-content: center;
  align-items: center;
  color: ${(p) => p.color && p.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
  width: 100%;
`;
