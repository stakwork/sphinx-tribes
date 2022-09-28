import React, { useEffect } from 'react';
import styled from 'styled-components';
import { EuiSwitch } from '@elastic/eui';
import type { Props } from './propsType';
import { FieldEnv, Note } from './index';
import { wantedCodingTaskSchema } from '../schema';

export default function SwitchInput({
  label,
  note,
  value,
  name,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  extraHTML
}: Props) {
  useEffect(() => {
    // if value not initiated, default value true
    if (name === 'show' && value === undefined) handleChange(true);

    // if (name === 'github_description') {
    //   wantedCodingTaskSchema.map((val) => {
    //     if (val.name === 'description') {
    //       console.log(val.name, value);
    //       return { ...val, type: value ? 'hide' : 'textarea' };
    //     }
    //   });
    // }
  }, []);

  return (
    <>
      <FieldEnv label={label}>
        <div style={{ padding: 10 }}>
          <EuiSwitch
            label=""
            checked={value}
            onChange={(e) => {
              handleChange(e.target.checked);
            }}
            onBlur={handleBlur}
            onFocus={handleFocus}
            compressed
            style={{
              border: 'none',
              background: 'inherit'
            }}
          />
        </div>
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
