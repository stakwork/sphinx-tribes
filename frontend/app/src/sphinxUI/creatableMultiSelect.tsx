import React from 'react';
import styled from 'styled-components';
import CreatableSelect from 'react-select/creatable';
import { colors } from '../colors';
import { StylesConfig } from 'react-select';
import { ColourOption, colourOptions } from '../people/utils/language_label_style';
// import makeAnimated from 'react-select/animated';

export default function Sel(props: any) {
  const { options, onChange, value, style, setIsTop } = props;
  const color = colors['light'];

  const opts =
    colourOptions.map((o) => {
      return {
        value: o.value,
        label: o.label,
        color: o.color,
        background: o.background,
        border: o.border
      };
    }) || [];

  return (
    <div style={{ position: 'relative', ...style }}>
      <S
        color={color}
        closeMenuOnSelect={false}
        isMulti
        options={opts}
        value={value}
        onChange={(value) => onChange(value)}
        onBlur={() => {
          setIsTop(false);
        }}
        onFocus={() => {
          setIsTop(true);
        }}
        className={'multi-select-input'}
        styles={{
          control: (styles) => ({ ...styles, backgroundColor: 'white' }),
          option: (styles, { data }: any) => {
            return {
              ...styles,
              backgroundColor: '#fff',
              color: color.text2,
              fontFamily: 'Barlow',
              fontSize: '14px',
              fontWeight: '500',
              ':hover': {
                background: color.light_blue200
              }
            };
          },
          multiValue: (styles, { data }: any) => {
            return {
              ...styles,
              backgroundColor: data.background,
              border: data.border,
              color: data.color,
              fontFamily: 'Barlow',
              fontSize: '14px',
              fontWeight: '500'
            };
          },
          multiValueLabel: (styles, { data }: any) => ({
            ...styles,
            background: data.background,
            color: data.color,
            border: data.border
          }),
          multiValueRemove: (styles, { data }: any) => ({
            ...styles,
            color: data.color,
            backgroundColor: data.background,
            ':hover': {
              backgroundColor: data.background,
              color: data.color
            }
          })
        }}
      />
    </div>
  );
}

interface styledProps {
  color?: any;
}

const S = styled(CreatableSelect)<styledProps>`
background:#ffffff00;
border: 1px solid ${(p) => p.color && p.color.grayish.G750};
color: ${(p) => p.color && p.color.pureBlack};
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
    background: ${(p) => p.color && p.color.pureWhite} !important;
    background-color: ${(p) => p.color && p.color.pureWhite} !important;
}
}
`;
