import React, { Fragment } from 'react';
import styled from 'styled-components';
import { EuiSuperSelect, EuiText } from '@elastic/eui';

export default function Select(props: any) {
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
                color: '#000',
                fontSize: '14px',
                padding: '0px',
                margin: 0
              }}>
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
                }}>
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

//euiContextMenuItem euiSuperSelect__item euiSuperSelect__item--hasDividers

const S = styled(EuiSuperSelect as any)`
  background: #ffffff00;
  border: 1px solid #e0e0e0;
  color: #000;
  box-sizing: border-box;
  box-shadow: none;

  user-select: none;

  .euiSuperSelectControl.euiSuperSelect--isOpen__button {
    background: #ffffff !important;
    background-color: #ffffff !important;
  }
  .euiPanel {
    background: #fff;
  }
  button {
    background: #ffffff !important;
    background-color: #ffffff !important;
  }
`;
