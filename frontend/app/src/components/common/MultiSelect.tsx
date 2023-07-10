import React from 'react';
import styled from 'styled-components';
import Select from 'react-select';
import { SelProps } from 'components/interfaces';
import { colors } from '../../config/colors';

interface styledProps {
  color?: any;
}

const S = styled(Select)<styledProps>`
background:#ffffff00;
border: 1px solid ${(p: any) => p.color && p.color.grayish.G750};
color: ${(p: any) => p.color && p.color.pureBlack};
box-sizing: border-box;
box-shadow:none;
border: none !important;
user-select:none;
font-size:12px;
border-width:0px !important;

#react-select-10-listbox{
    z-index:1000;
}
#react-select-9-listbox{
    z-index:1000;
}
#react-select-8-listbox{
    z-index:1000;
}
#react-select-7-listbox{
    z-index:1000;
}
#react-select-6-listbox{
    z-index:1000;
}
#react-select-5-listbox{
    z-index:1000;
}
#react-select-4-listbox{
    z-index:1000;
}
#react-select-3-listbox{
    z-index:1000;
}
#react-select-2-listbox{
    z-index:1000;
}
#react-select-1-listbox{
    z-index:1000;
}

div {
    border-width:0px !important;
    border: none !important;
}

button {
    background: ${(p: any) => p.color && p.color.pureWhite} !important;
    background-color: ${(p: any) => p.color && p.color.pureWhite} !important;
}
}
`;
export default function Sel(props: SelProps) {
  const { options, onChange, value, style } = props;
  const color = colors['light'];

  const opts = options
    ? options.map((o: any) => ({
        value: o.value,
        label: o.label
      }))
    : [];

  return (
    <div style={{ position: 'relative', ...style }}>
      <S
        color={color}
        closeMenuOnSelect={false}
        isMulti
        options={opts}
        value={value}
        onChange={(value: any) => onChange(value)}
        className={'multi-select-input'}
      />
    </div>
  );
}
