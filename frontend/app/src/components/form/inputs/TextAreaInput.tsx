import React, { useState } from 'react';
import styled from 'styled-components';
import { EuiIcon } from '@elastic/eui';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { FieldEnv, FieldTextArea, Note } from './index';

const StyleOnText = {
  'Description *': {
    height: '172px',
    width: '292px'
  },
  Deliverables: {
    height: '135px',
    width: '192px'
  }
};

const defaultHeight = '135px';
const defaultWidth = '100%';

interface styledProps {
  color?: any;
}

const OuterContainer = styled.div<styledProps>`
  .euiFormRow_active {
    border: 1px solid ${(p: any) => p?.color && p?.color.blue2};
    .euiFormRow__labelWrapper {
      margin-bottom: 0px;
      margin-top: -9px;
      padding-left: 10px;
      height: 14px;
      label {
        color: ${(p: any) => p?.color && p?.color.grayish.G300} !important;
        background: ${(p: any) => p?.color && p?.color.pureWhite};
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
        color: ${(p: any) => p?.color && p?.color.grayish.G300} !important;
        background: ${(p: any) => p?.color && p?.color.pureWhite};
        z-index: 10;
      }
    }
  }
`;

const ExtraText = styled.div<styledProps>`
  color: ${(p: any) => p?.color && p?.color.grayish.G760};
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
  color: ${(p: any) => p?.color && p?.color.blue3};
  pointer-events: none;
  user-select: none;
`;
const R = styled.div`
  position: relative;
`;
export default function TextAreaInput({
  error,
  note,
  label,
  value,
  handleChange,
  handleBlur,
  handleFocus,
  readOnly,
  extraHTML,
  borderType
}: Props) {
  let labeltext = label;
  const color = colors['light'];
  if (error) labeltext = `${labeltext} (${error})`;
  const [active, setActive] = useState<boolean>(false);
  return (
    <OuterContainer color={color}>
      <FieldEnv
        color={color}
        onClick={() => {
          setActive(true);
        }}
        className={active ? 'euiFormRow_active' : (value ?? '') === '' ? '' : 'euiFormRow_filed'}
        border={borderType}
        label={label}
        height={StyleOnText[label]?.height ?? defaultHeight}
        width={StyleOnText[label]?.width ?? defaultWidth}
      >
        <R>
          <FieldTextArea
            color={color}
            height={StyleOnText[label]?.height ?? defaultHeight}
            width={StyleOnText[label]?.width ?? defaultWidth}
            name="first"
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
            rows={label === 'Description' ? 8 : 6}
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
        color={color}
        style={{ display: value && extraHTML ? 'block' : 'none' }}
        dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
      />
    </OuterContainer>
  );
}
