import React from 'react';
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

  const padStyle = prepend ? { paddingLeft: 0 } : {};
  return (
    <>
      <FieldEnv border={borderType} label={labeltext}>
        <R>
          <FieldText
            name="first"
            value={value || ''}
            readOnly={readOnly || false}
            onChange={(e) => handleChange(e.target.value)}
            onBlur={handleBlur}
            onFocus={handleFocus}
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
    </>
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
