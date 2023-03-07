import React, { useState } from 'react';
import styled from 'styled-components';
import type { Props } from './propsType';
import { FieldEnv, FieldText, Note } from './index';
import { satToUsd } from '../../../helpers';
import { colors } from '../../../config/colors';

export default function NumberInput({
  name,
  error,
  note,
  label,
  value,
  extraHTML,
  handleChange,
  handleBlur,
  handleFocus,
  borderType
}: Props) {
  let labeltext = label;
  if (error) labeltext = `${labeltext} (${error})`;
  const [active, setActive] = useState<boolean>(false);
  const color = colors['light'];

  return (
    <OuterContainer color={color}>
      <FieldEnv
        color={color}
        onClick={() => {
          setActive(true);
        }}
        className={active ? 'euiFormRow_active' : (value ?? '') === '' ? '' : 'euiFormRow_filed'}
        border={borderType}
        label={labeltext}
      >
        <R>
          <FieldText
            color={color}
            name="first"
            value={value}
            type="number"
            onChange={(e) => {
              // dont allow zero or negative numbers
              if (parseInt(e.target.value) < 0) return;
              handleChange(e.target.value);
            }}
            onBlur={(e) => {
              // enter 0 on blur if no value
              console.log('onBlur', value);
              if (value === '') handleChange(0);
              if (value === '0') handleChange(0);
              handleBlur(e);
              setActive(false);
            }}
            onFocus={(e) => {
              // remove 0 on focus
              console.log('onFocus', value);
              if (value === 0) handleChange('');
              handleFocus(e);
              setActive(true);
            }}
          />
        </R>
      </FieldEnv>
      {note && <Note color={color}>*{note}</Note>}
      {name.includes('price') && <Note color={color}>({satToUsd(value)} USD)</Note>}
      <ExtraText
        color={color}
        style={{ display: value && extraHTML ? 'block' : 'none' }}
        dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
      />
    </OuterContainer>
  );
}

interface styledProps {
  color?: any;
}

const OuterContainer = styled.div<styledProps>`
  .euiFormRow_active {
    border: 1px solid ${(p) => p.color && p.color.blue2};
    .euiFormRow__labelWrapper {
      margin-bottom: 0px;
      margin-top: -9px;
      padding-left: 10px;
      height: 14px;
      label {
        color: ${(p) => p.color && p.color.grayish.G300} !important;
        background: ${(p) => p.color && p.color.pureWhite};
        z-index: 10;
      }
    }
  }
  .euiFormRow_filed {
    .euiFormRow__labelWrapper {
      margin-bottom: 0px;
      margin-top: -9px;
      padding-left: 10px;
      height: 14px;
      label {
        color: ${(p) => p.color && p.color.grayish.G300} !important;
        background: ${(p) => p.color && p.color.pureWhite};
        z-index: 10;
      }
    }
  }
`;

const ExtraText = styled.div<styledProps>`
  padding: 0px 10px 5px;
  margin: -5px 0 10px;
  color: ${(p) => p.color && p.color.red3};
  font-style: italic;
  max-width: calc(100% - 20px);
  word-break: break;
  font-size: 14px;
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
`;
