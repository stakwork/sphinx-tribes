import React from 'react';
import styled from 'styled-components';
import Select from 'react-select';
import { SelProps } from 'components/interfaces';

const S = styled(Select)`
background:#ffffff00;
border: 1px solid #E0E0E0;
color:#000;
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
    background:#ffffff !important;
    background-color:#ffffff !important;
}
}
`;
export default function SearchableSelect(props: SelProps) {
  const { options, onChange, onInputChange, value, style, loading } = props;

  const opts = options
    ? options.map((o: any) => ({
        ...o,
        value: o.value,
        label: o.label
      }))
    : [];

  return (
    <div style={{ position: 'relative', ...style }}>
      <S
        options={opts}
        isLoading={loading}
        placeholder={'Type to search...'}
        isClearable={true}
        isSearchable={true}
        value={value}
        onChange={onChange}
        onInputChange={onInputChange}
        className={'searchable-select-input'}
      />
    </div>
  );
}
