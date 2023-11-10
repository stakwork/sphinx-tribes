import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiIcon } from '@elastic/eui';
import { CreatableMultiSelect } from '../../common';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { FieldEnv, Note } from './index';

interface styledProps {
  color?: any;
}

const ExtraText = styled.div`
  padding: 2px 10px 25px 10px;
  max-width: calc(100 % - 20px);
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
  color: ${(p: any) => p?.color && p.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
`;
export default function CreatableMultiSelectInput({
  error,
  note,
  type,
  label,
  options,
  value,
  handleChange,
  extraHTML
}: Props) {
  let labeltext = label;
  if (error) labeltext = `${labeltext} (INCORRECT FORMAT)`;
  const color = colors['light'];

  const [isTop, setIsTop] = useState<boolean>(false);

  return (
    <>
      <FieldEnv
        color={color}
        label={labeltext}
        isTop={isTop || value?.length > 0}
        isFill={value?.length > 0}
      >
        <R>
          <CreatableMultiSelect
            selectStyle={{ border: 'none' }}
            options={options}
            writeMode={type === 'multiselectwrite'}
            value={value}
            onChange={(e: any) => {
              handleChange(e);
              setIsTop(true);
            }}
            setIsTop={setIsTop}
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
