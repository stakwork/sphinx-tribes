import React, { useEffect } from 'react';
import styled from 'styled-components';
import { EuiSwitch, EuiText } from '@elastic/eui';
import { colors } from '../../../config/colors';
import type { Props } from './propsType';
import { Note } from './index';

interface styledProps {
  color?: any;
}

const ExtraText = styled.div<styledProps>`
  color: ${(p: any) => p?.color && p?.color.grayish.G760};
  padding: 10px 10px 25px 10px;
  max-width: calc(100% - 20px);
  word-break: break-all;
  font-size: 14px;
`;

const Container = styled.div<styledProps>`
  padding: 10px;
  display: flex;
  align-items: center;
  .Label {
    font-family: 'Barlow';
    font-style: normal;
    font-weight: 500;
    font-size: 14px;
    line-height: 35px;
    display: flex;
    align-items: center;
    color: #292c33;
    margin-right: 4px;
  }
`;
export default function SwitchInput({
  label,
  note,
  value,
  name,
  handleChange,
  handleBlur,
  handleFocus,
  extraHTML,
  disabled,
  style = {}
}: Props) {
  useEffect(() => {
    // if value not initiated, default value true
    if (name === 'show' && value === undefined) handleChange(true);
  }, [handleChange, name, value]);

  const color = colors['light'];

  return (
    <>
      <Container style={style} color={color}>
        <EuiText className="Label">{label}</EuiText>
        <EuiSwitch
          className="switcher"
          label=""
          checked={value}
          onChange={(e: any) => {
            handleChange(e.target.checked);
          }}
          onBlur={handleBlur}
          onFocus={handleFocus}
          compressed
          disabled={disabled}
          style={{
            border: 'none',
            background: 'inherit'
          }}
        />
      </Container>
      {note && <Note color={color}>*{note}</Note>}
      <ExtraText
        style={{ display: value && extraHTML ? 'block' : 'none' }}
        dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
      />
    </>
  );
}
