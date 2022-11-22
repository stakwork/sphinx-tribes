import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiTextArea, EuiIcon } from '@elastic/eui';
import type { Props } from './propsType';
import { FieldEnv, FieldTextArea, Note } from './index';
import { colors } from '../../colors';

export default function TextAreaInput({
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
  // console.log("TEXTAREA", label, extraHTML)

  let labeltext = label;
  const color = colors['light'];
  if (error) labeltext = `${labeltext} (${error})`;
  const [active, setActive] = useState<boolean>(false);
  return (
    <OuterContainer>
      <FieldEnv
        onClick={() => {
          setActive(true);
        }}
        className={active ? 'euiFormRow_active' : (value ?? '') === '' ? '' : 'euiFormRow_filed'}
        border={borderType}
        label={labeltext}>
        <R>
          <FieldTextArea
            name="first"
            value={value || ''}
            readOnly={readOnly || false}
            onChange={(e) => handleChange(e.target.value)}
            onBlur={(e) => {
              handleBlur(e);
              setActive(false);
            }}
            onFocus={(e) => {
              handleFocus(e);
              setActive(true);
            }}
            rows={4}
            // prepend={prepend}
          />
          {error && (
            <E>
              <EuiIcon type="alert" size="m" style={{ width: 20, height: 20 }} />
            </E>
          )}
        </R>
      </FieldEnv>
      {note && <Note>*{note}</Note>}
      <ExtraText
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
    border: 1px solid ${(p) => p?.color && p?.color.blue2};
    .euiFormRow__labelWrapper {
      margin-bottom: 0px;
      margin-top: -9px;
      padding-left: 10px;
      height: 14px;
      label {
        color: ${(p) => p?.color && p?.color.grayish.G300} !important;
        background: ${(p) => p?.color && p?.color.pureWhite};
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
        color: ${(p) => p?.color && p?.color.grayish.G300} !important;
        background: ${(p) => p?.color && p?.color.pureWhite};
        z-index: 10;
      }
    }
  }
`;

const ExtraText = styled.div<styledProps>`
  color: ${(p) => p?.color && p?.color.grayish.G760};
  padding: 10px 10px 25px 10px;
  max-width: calc(100% - 20px);
  word-break: break-all;
  font-size: 14px;
`;

const E = styled.div<styledProps>`
  position: absolute;
  right: 10px;
  top: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: ${(p) => p?.color && p?.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
`;
