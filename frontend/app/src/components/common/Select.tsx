import React from 'react';
import styled from 'styled-components';
import { EuiSuperSelect, EuiText } from '@elastic/eui';
import { SelProps } from 'components/interfaces';
import { colors } from '../../config/colors';

interface styleProps {
  color?: any;
}

const S = styled(EuiSuperSelect as any)<styleProps>`
  background: ${(p: any) => p?.color && p.color.pureWhite};
  border: 1px solid ${(p: any) => p?.color && p?.color.grayish.G750};
  color: ${(p: any) => p?.color && p?.color.pureBlack};
  box-sizing: border-box;
  box-shadow: none;
  padding-left: 16px;
  user-select: none;
  .euiSuperSelectControl.euiSuperSelect--isOpen__button {
    background: ${(p: any) => p?.color && p?.color.pureWhite} !important;
    background-color: ${(p: any) => p?.color && p?.color.pureWhite} !important;
  }
  .euiPanel {
    background: ${(p: any) => p?.color && p?.color.pureWhite};
  }
  .button {
    background: ${(p: any) => p?.color && p?.color.pureWhite} !important;
    background-color: ${(p: any) => p?.color && p?.color.pureWhite} !important;
  }
`;
export default function Select(props: SelProps) {
  const color = colors['light'];
  const { options, onChange, value, style, selectStyle, handleActive, testId } = props;

  const opts = options
    ? options.map((o: any) => ({
        value: o.value,
        inputDisplay: o.label,
        dropdownDisplay: (
          <>
            <p
              style={{
                color: color.text2,
                fontSize: '14px',
                paddingLeft: '0px',
                margin: 0,
                fontFamily: 'Barlow',
                fontWeight: '500',
                lineHeight: '32px',
                letterSpacing: '0.01em'
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
      }))
    : [];

  return (
    <div style={{ position: 'relative', ...style }}>
      <S
        data-testid={testId}
        color={color}
        style={{
          ...selectStyle
        }}
        onFocus={() => {
          if (handleActive) handleActive(true);
        }}
        onBlur={() => {
          if (handleActive) handleActive(false);
        }}
        options={opts}
        valueOfSelected={value}
        onChange={(value: any) => {
          onChange(value);
          if (handleActive) handleActive(false);
        }}
        fullWidth={true}
        itemLayoutAlign="top"
      />
    </div>
  );
}
