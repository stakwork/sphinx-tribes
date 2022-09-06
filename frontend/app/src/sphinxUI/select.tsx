import React, { Fragment } from "react";
import styled from "styled-components";
import { EuiSuperSelect, EuiText } from "@elastic/eui";

export default function Select(props: any) {
  const { options, onChange, value, style, selectStyle } = props;

  const opts =
    options.map((o) => {
      return {
        value: o.value,
        inputDisplay: o.label,
        dropdownDisplay: (
          <Fragment>
            <strong>{o.label}</strong>
            <EuiText size="s" color="subdued">
              <p className="euiTextColor--subdued">{o.description}</p>
            </EuiText>
          </Fragment>
        ),
      };
    }) || [];

  return (
    <div style={{ position: "relative", ...style }}>
      <S
        style={selectStyle}
        options={opts}
        valueOfSelected={value}
        onChange={(value) => onChange(value)}
        itemLayoutAlign="top"
        hasDividers
      />
    </div>
  );
}

const S = styled(EuiSuperSelect as any)`
background:#ffffff00;
border: 1px solid #E0E0E0;
color:#000;
box-sizing: border-box;
box-shadow:none;

user-select:none;

.euiSuperSelectControl.euiSuperSelect--isOpen__button{
    background:#ffffff !important;
    background-color:#ffffff !important;
}

button {
    background:#ffffff !important;
    background-color:#ffffff !important;
}
}
`;
