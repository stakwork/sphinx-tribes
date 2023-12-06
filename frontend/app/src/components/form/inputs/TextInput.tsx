import React, { useState } from 'react';
import styled from 'styled-components';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { FieldEnv, FieldText, Note } from './index';

interface styledProps {
  color?: any;
}

const OuterContainer = styled.div<styledProps>`
  .euiFormRow_active {
    border: 1px solid ${(p: any) => p.color && p.color.blue2};
    .euiFormRow__labelWrapper {
      margin-bottom: 0px;
      margin-top: -9px;
      padding-left: 10px;
      height: 14px;
      label {
        color: ${(p: any) => p.color && p.color.grayish.G300} !important;
        background: ${(p: any) => p.color && p.color.pureWhite};
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
        color: ${(p: any) => p.color && p.color.grayish.G300} !important;
        background: ${(p: any) => p.color && p.color.pureWhite};
        z-index: 10;
      }
    }
  }
`;

const ExtraText = styled.div<styledProps>`
  padding: 0px 10px 5px;
  margin: -5px 0 10px;
  color: ${(p: any) => p.color && p.color.red3};
  font-style: italic;
  max-width: calc(100% - 20px);
  word-break: break-all;
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
  color: ${(p: any) => p.color && p.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
`;
export default function TextInput({
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
  const [active, setActive] = useState<boolean>(false);
  const color = colors['light'];

  const padStyle = prepend ? { paddingLeft: 0 } : {};
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
        isTextField={true}
        error={error}
      >
        <R>
          <FieldText
            color={color}
            name={'first'}
            value={value || ''}
            readOnly={readOnly || false}
            onChange={(e: any) => handleChange(e.target.value)}
            onBlur={(e: any) => {
              handleBlur(e);
              setActive(false);
            }}
            onFocus={(e: any) => {
              handleFocus(e);
              setActive(true);
            }}
            prepend={active ? prepend : ''}
            style={padStyle}
            isTextField={true}
          />
          {error && <E color={color} />}
        </R>
      </FieldEnv>
      {note && <Note color={color}>*{note}</Note>}
      <ExtraText
        color={color}
        style={{ display: value && extraHTML ? 'block' : 'none' }}
        dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
      />
    </OuterContainer>
  );
}
