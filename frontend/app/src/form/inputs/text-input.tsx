import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiFieldText, EuiIcon } from '@elastic/eui';
import type { Props } from './propsType';
import { FieldEnv, FieldText, Note } from './index';

export default function TextInput({
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
  const [active, setActive] = useState<boolean>(false);

  const padStyle = prepend ? { paddingLeft: 0 } : {};
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
          <FieldText
            name={'first'}
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
            prepend={prepend}
            style={padStyle}
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

const OuterContainer = styled.div`
  .euiFormRow_active {
    border: 1px solid #82b4ff;
    .euiFormRow__labelWrapper {
      margin-bottom: 0px;
      margin-top: -9px;
      padding-left: 10px;
      height: 14px;
      label {
        color: #b0b7bc !important;
        background: #ffffff;
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
        color: #b0b7bc !important;
        background: #ffffff;
        z-index: 10;
      }
    }
  }
`;

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
