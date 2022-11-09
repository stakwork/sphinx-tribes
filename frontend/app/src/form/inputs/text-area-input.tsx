import React from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiTextArea, EuiIcon } from '@elastic/eui';
import type { Props } from './propsType';
import { FieldEnv, FieldTextArea, Note } from './index';

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
  if (error) labeltext = `${labeltext} (${error})`;

  return (
    <>
      <FieldEnv border={borderType} label={labeltext}>
        <R>
          <FieldTextArea
            name="first"
            value={value || ''}
            readOnly={readOnly || false}
            onChange={(e) => handleChange(e.target.value)}
            onBlur={handleBlur}
            onFocus={handleFocus}
            rows={2}
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
    </>
  );
}

const ExtraText = styled.div`
  color: #ddd;
  padding: 10px 10px 25px 10px;
  max-width: calc(100% - 20px);
  word-break: break-all;
  font-size: 14px;
`;

const E = styled.div`
  position: absolute;
  right: 10px;
  top: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #45b9f6;
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
`;
