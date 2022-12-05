import React from 'react';
import styled from 'styled-components';
import { EuiIcon } from '@elastic/eui';
import type { Props } from './propsType';
import { FieldEnv, Note } from './index';
import { MultiSelect } from '../../sphinxUI';
import { colors } from '../../colors';

export default function MultiSelectInput({
  error,
  note,
  type,
  label,
  options,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  prepend,
  extraHTML
}: Props) {
  let labeltext = label;
  if (error) labeltext = `${labeltext} (${error})`;
  const color = colors['light'];

  return (
    <>
      <FieldEnv label={labeltext} color={color}>
        <R>
          <MultiSelect
            selectStyle={{ border: 'none' }}
            options={options}
            writeMode={type === 'multiselectwrite'}
            value={value}
            onChange={(e) => {
              console.log('onChange', e);
              handleChange(e);
            }}
          />
          {error && (
            <E color={color}>
              <EuiIcon type="alert" size="m" style={{ width: 20, height: 20 }} />
            </E>
          )}
        </R>
      </FieldEnv>
      {note && <Note color={color}>*{note}</Note>}
      <ExtraText
        style={{ display: value && extraHTML ? 'block' : 'none' }}
        dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
      />
    </>
  );
}

interface styledProps {
  color?: any;
}

const ExtraText = styled.div`
  padding: 2px 10px 25px 10px;
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
  color: ${(p) => p.color && p.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
`;
