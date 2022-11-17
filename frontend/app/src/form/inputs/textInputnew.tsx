import React from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiFieldText, EuiIcon } from '@elastic/eui';
import { Props } from './propsType';
import { FieldEnv, FieldText, Note } from '.';

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
  if (error) labeltext = `${labeltext} (${error})`;

  const padStyle = prepend ? { paddingLeft: 0 } : {};
  return (
    <TextContainer>
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

const ExtraText = styled.div`
  padding: 0px 10px 5px;
  margin: -5px 0 10px;
  color: #b75858;
  font-style: italic;
  max-width: calc(100% - 20px);
  word-break: break-all;
  font-size: 14px;
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
`;

const TextContainer = styled.div`
  margin-bottom: 13px;
  .label-float {
    position: relative;
    padding-top: 13px;
    font-size: 13px;
    color: #b0b7bc;
    font-family: Barlow;
  }

  .label-float input {
    border: 1px solid #dde1e5;
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
    border: 1px solid #82b4ff;
  }

  .label-float input::placeholder {
    color: #b0b7bc;
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
    color: #b0b7bc;
  }
`;
