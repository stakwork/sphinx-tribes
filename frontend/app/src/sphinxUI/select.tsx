import React, { Fragment } from 'react';
import styled from 'styled-components';
import { EuiSuperSelect, EuiText } from '@elastic/eui';
import { colors } from '../colors';

export default function Select(props: any) {
  const color = colors['light'];
  const { options, onChange, value, style, selectStyle, handleActive } = props;

  const opts =
    options.map((o) => {
      return {
        value: o.value,
        inputDisplay: o.label,
        dropdownDisplay: (
          <>
            <p
              style={{
                color: color.pureBlack,
                fontSize: '14px',
                padding: '0px',
                margin: 0
              }}
            >
              {o.label}
            </p>
            {o.description && (
              <EuiText
                size="s"
                color="subdued"
                style={{
                  padding: 0,
                  margin: 0,
                  fontSize: '12px'
                }}
              >
                <p className="euiTextColor--subdued">{o.description}</p>
              </EuiText>
            )}
          </>
        )
      };
    }) || [];

  return (
    <div style={{ position: 'relative', ...style }}>
      <S
        color={color}
        style={{
          ...selectStyle
        }}
        onFocus={() => {
          handleActive(true);
        }}
        onBlur={() => {
          handleActive(false);
        }}
        options={opts}
        valueOfSelected={value}
        onChange={(value) => {
          onChange(value);
          handleActive(false);
        }}
        itemLayoutAlign="top"
      />
    </div>
  );
}

interface styleProps {
  color?: any;
}

const S = styled(EuiSuperSelect as any)<styleProps>`
  background: ${(p) => p?.color && p.color.pureWhite};
  border: 1px solid ${(p) => p?.color && p?.color.grayish.G750};
  color: ${(p) => p?.color && p?.color.pureBlack};
  box-sizing: border-box;
  box-shadow: none;
  user-select: none;
  .euiSuperSelectControl.euiSuperSelect--isOpen__button {
    background: ${(p) => p?.color && p?.color.pureWhite} !important;
    background-color: ${(p) => p?.color && p?.color.pureWhite} !important;
  }
  .euiPanel {
    background: ${(p) => p?.color && p?.color.pureWhite};
  }
  button {
    background: ${(p) => p?.color && p?.color.pureWhite} !important;
    background-color: ${(p) => p?.color && p?.color.pureWhite} !important;
  }
`;
