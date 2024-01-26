import React from 'react';
import styled from 'styled-components';
import { colors } from '../../../config/colors';
import { Props } from './propsType';

interface styledProps {
  color?: any;
}
const TextContainer = styled.div<styledProps>`
  margin-bottom: 13px;
  .label-float {
    position: relative;
    padding-top: 13px;
    font-size: 13px;
    color: ${(p: any) => p?.color && p?.color.grayish.G300};
    font-family: 'Barlow';
  }

  .label-float input {
    border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G600};
    border-radius: 4px;
    outline: none;
    min-width: 290px;
    padding: 9.5px 20px;
    font-size: 13px;
    font-family: 'Barlow';
    transition: all 0.1s linear;
    -webkit-transition: all 0.1s linear;
    -moz-transition: all 0.1s linear;
    -webkit-appearance: none;
  }

  .label-float input:focus {
    border: 1px solid ${(p: any) => p?.color && p?.color.blue2};
  }

  .label-float input::placeholder {
    color: ${(p: any) => p?.color && p?.color.grayish.G300};
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
    color: ${(p: any) => p?.color && p?.color.grayish.G300};
  }
`;
export default function TextInputNew({ label, value, handleChange, readOnly }: Props) {
  const color = colors['light'];
  return (
    <TextContainer color={color}>
      <div className="label-float">
        <input
          type="text"
          placeholder=" "
          value={value || ''}
          readOnly={readOnly || false}
          onChange={(e: any) => handleChange(e.target.value)}
        />
        <label>{label}</label>
      </div>
    </TextContainer>
  );
}
