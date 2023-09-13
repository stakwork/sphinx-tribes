import MaterialIcon from '@material/react-material-icon';
import React, { useState } from 'react';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { SearchTextInputProps } from 'components/interfaces';
import { useStores } from '../../store';
let debounceValue = '';
let inDebounce;
function debounce(func: any, delay: any) {
  clearTimeout(inDebounce);
  inDebounce = setTimeout(() => {
    func();
  }, delay);
}

const Text = styled.input`
  background: #f2f3f580;
  border: 1px solid #e0e0e0;
  box-sizing: border-box;
  border-radius: 21px;
  padding-left: 20px;
  padding-right: 10px;
  // width:100%;
  font-style: normal;
  font-weight: normal;
  font-size: 12px;
  line-height: 14px;
  height: 35px;
  transition: all 0.4s;
`;

function SearchTextInput(props: SearchTextInputProps) {
  const { ui } = useStores();
  const [searchValue, setSearchValue] = useState(ui.searchText || '');
  const [expand, setExpand] = useState(ui.searchText ? true : false);

  const collapseStyles = props.small && !expand ? {} : {};

  function doDelayedValueUpdate() {
    props.onChange(debounceValue);
  }

  function erase() {
    setSearchValue('');
    props.onChange('');
  }

  return (
    <div style={{ position: 'relative' }}>
      <Text
        {...props}
        onFocus={() => setExpand(true)}
        onBlur={() => {
          if (!ui.searchText) setExpand(false);
        }}
        value={searchValue}
        onChange={(e: any) => {
          setSearchValue(e.target.value);
          debounceValue = e.target.value;
          debounce(doDelayedValueUpdate, 300);
        }}
        placeholder={'Search'}
        style={{ ...props.style, ...collapseStyles }}
      />
      {searchValue ? (
        <MaterialIcon
          icon="close"
          onClick={() => erase()}
          style={{
            position: 'absolute',
            color: '#757778',
            cursor: 'pointer',
            top: 9,
            right: 9,
            fontSize: 22,
            userSelect: 'none',
            ...props.iconStyle
          }}
        />
      ) : (
        <MaterialIcon
          icon="search"
          style={{
            position: 'absolute',
            color: '#B0B7BC',
            top: 9,
            right: 9,
            fontSize: 22,
            userSelect: 'none',
            pointerEvents: 'none',
            ...props.iconStyle
          }}
        />
      )}
    </div>
  );
}

export default observer(SearchTextInput);
