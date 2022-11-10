import React from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiFieldText, EuiIcon, PropertySortType } from '@elastic/eui';
import type { Props } from './propsType';

export default function SearchTextInput({
  placeholder,
  error,
  label,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  extraHTML
}: any) {
  let labeltext = label;
  if (error) labeltext = `${labeltext} (${error})`;

  return (
    <>
      <R>
        <Text
          name="first"
          value={value || ''}
          readOnly={readOnly || false}
          onChange={(e) => handleChange(e.target.value)}
          onBlur={handleBlur}
          onFocus={handleFocus}
          placeholder={placeholder || 'Search'}
        />
        {error && (
          <E>
            <EuiIcon type="alert" size="m" style={{ width: 20, height: 20 }} />
          </E>
        )}
      </R>
    </>
  );
}

const Text = styled.input`
  background: #f2f3f580;
  border: 1px solid #e0e0e0;
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

const E = styled.div`
  position: absolute;
  right: 10px;
  top: 0px;
  display: flex;
  height: 100%;
  justify-content: center;
  align-items: center;
  color: #45b9f6;
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
  width: 100%;
`;
