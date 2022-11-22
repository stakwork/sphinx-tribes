import React from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiFieldText, EuiIcon } from '@elastic/eui';
import { Props } from './propsType';
import { FieldEnv, FieldText, Note } from '.';
import { colors } from '../../colors';

export default function TextInputNew({
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
  extraHTML,
  borderType
}: Props) {
  let labeltext = label;
  const color = colors['light'];
  if (error) labeltext = `${labeltext} (${error})`;

  const padStyle = prepend ? { paddingLeft: 0 } : {};
  return (
    <TextContainer color={color}>
      <div className="label-float">
        <input
          type="text"
          placeholder=" "
          value={value || ''}
          readOnly={readOnly || false}
          onChange={(e) => handleChange(e.target.value)}
        />
        <label>{label}</label>
      </div>
    </TextContainer>
  );
}

interface styledProps {
  color?: any;
}

const TextContainer = styled.div<styledProps>`
  margin-bottom: 13px;
  .label-float {
    position: relative;
    padding-top: 13px;
    font-size: 13px;
    color: ${(p) => p?.color && p?.color.grayish.G300};
    font-family: Barlow;
  }

  .label-float input {
    border: 1px solid ${(p) => p?.color && p?.color.grayish.G600};
    border-radius: 4px;
    outline: none;
    min-width: 290px;
    padding: 9.5px 20px;
    font-size: 13px;
    font-family: Barlow;
    transition: all 0.1s linear;
    -webkit-transition: all 0.1s linear;
    -moz-transition: all 0.1s linear;
    -webkit-appearance: none;
  }

  .label-float input:focus {
    border: 1px solid ${(p) => p?.color && p?.color.blue2};
  }

  .label-float input::placeholder {
    color: ${(p) => p?.color && p?.color.grayish.G300};
  }

  .label-float label {
    pointer-events: none;
    position: absolute;
    top: calc(50% - 8px);
    left: 15px;
    transition: all 0.1s linear;
    -webkit-transition: all 0.1s linear;
    -moz-transition: all 0.1s linear;
    background-color: white;
    padding: 5px;
    box-sizing: border-box;
  }

  .label-float input:required:invalid + label {
    color: red;
  }
  .label-float input:focus:required:invalid {
    border: 2px solid red;
  }
  .label-float input:required:invalid + label:before {
    content: '*';
  }
  .label-float input:focus + label,
  .label-float input:not(:placeholder-shown) + label {
    font-size: 13px;
    top: 0;
    color: ${(p) => p?.color && p?.color.grayish.G300};
  }
`;
