import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { FieldEnv } from '.';

interface styledProps {
  color?: any;
  borderColor?: string;
}

const InputOuterBox = styled.div<styledProps>`
  position: relative;
  margin-bottom: 16px;
  .inputText {
    height: 45px;
    width: 292px;
    font-size: 14px;
    color: ${(p: any) => p.color && p.color.pureBlack};
    border: 1px solid ${(p: any) => p.borderColor && p.borderColor};
    border-radius: 4px;
    outline: none;
    padding-left: 16px;
    font-family: 'Barlow';
    font-weight: 500;
    color: #3c3f41;
    letter-spacing: 0.01em;

    :active {
      border: 1px solid ${(p: any) => p.color && p.color.blue2} !important;
    }
    :focus-visible {
      border: 1px solid ${(p: any) => p.color && p.color.blue2} !important;
    }
  }
  .USD {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 35px;
    display: flex;
    align-items: center;
    color: ${(p: any) => p.color && p.color.grayish.G300};
    border: 3px solid ${(p: any) => p.color && p.color.pureWhite};
    margin-left: 16px;
  }
`;
export default function NumberInputNew({
  error,
  label,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  borderType
}: Props) {
  let labeltext = label;
  if (error) labeltext = `${labeltext}*`;

  const color = colors['light'];
  const [isError, setIsError] = useState<boolean>(false);
  const [textValue, setTextValue] = useState(value);
  const [active, setActive] = useState<boolean>(false);

  useEffect(() => {
    if (textValue) {
      setIsError(false);
    }
  }, [textValue]);

  return (
    <InputOuterBox color={color} borderColor={isError ? color.red2 : color.grayish.G600}>
      <FieldEnv
        color={color}
        onClick={() => {
          setActive(true);
        }}
        className={active ? 'euiFormRow_active' : (value ?? '') === '' ? '' : 'euiFormRow_filed'}
        border={borderType}
        label={labeltext}
      >
        <input
          className="inputText"
          id={'number'}
          type={'number'}
          value={textValue}
          placeholder={'0'}
          onFocus={handleFocus}
          onBlur={() => {
            handleBlur();
            if (error) {
              setIsError(true);
            }
          }}
          onChange={(e: any) => {
            handleChange(e.target.value);
            setTextValue(e.target.value);
          }}
        />
      </FieldEnv>
    </InputOuterBox>
  );
}
