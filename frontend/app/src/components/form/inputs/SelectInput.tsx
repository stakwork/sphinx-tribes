import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiIcon } from '@elastic/eui';
import { Select } from '../../common';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { FieldEnv, Note } from './index';

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
  color: ${(p: any) => p?.color && p?.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  margin-top: 1px;
  position: relative;
`;

const OuterContainer = styled.div<styledProps>`
  box-shadow: 0px 1px 2px ${(p: any) => p.color && p.color.black100} !important ;

  .euiFormRow_filed {
    position: relative;
    .euiFormRow__labelWrapper {
      margin-bottom: 0px;
      margin-top: -10px;
      padding-left: 10px;
      height: 14px;
      transition: all 0.4s;
      label {
        color: ${(p: any) => p?.color && p?.color.grayish.G300} !important;
        background: ${(p: any) => p?.color && p?.color.pureWhite};
        z-index: 10;
        font-family: 'Barlow';
        font-size: 12px;
        font-weight: 500;
        margin-left: 6px;
      }
    }
  }
  .euiFormRow_active {
    padding: 1px 0;
    border: 1px solid ${(p: any) => p?.color && p?.color.blue2};
  }

  .euiFormControlLayoutCustomIcon {
    color: ${(p: any) => p.color && p.color.text2_4};
  }
`;
export default function SelectInput({
  error,
  note,
  label,
  options,
  value,
  handleChange,
  extraHTML,
  testId
}: Props) {
  let labeltext = label;
  const color = colors['light'];
  if (error) labeltext = `${labeltext} (${error})`;
  const [active, setActive] = useState<boolean>(false);
  return (
    <OuterContainer color={color}>
      <FieldEnv
        color={color}
        label={labeltext}
        className={value ? 'euiFormRow_filed' : active ? 'euiFormRow_active' : ''}
      >
        <R>
          <Select
            testId={testId}
            name={'first'}
            selectStyle={{
              border: 'none',
              fontFamily: 'Barlow',
              fontWeight: '500',
              fontSize: '14px',
              color: '#3C3F41',
              letterSpacing: '0.01em'
            }}
            options={options}
            value={value}
            handleActive={setActive}
            onChange={(e: any) => {
              handleChange(e);
              setActive(false);
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
    </OuterContainer>
  );
}
