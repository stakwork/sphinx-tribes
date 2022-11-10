import MaterialIcon from '@material/react-material-icon';
import React, { useState } from 'react';
import styled from 'styled-components';
import { useStores } from '../store';

export default function SearchTextInput(props: any) {
  const { ui } = useStores();
  const [searchValue, setSearchValue] = useState(ui.searchText || '');
  const [expand, setExpand] = useState(ui.searchText ? true : false);

  const collapseStyles =
    props.small && !expand
      ? {
          // width: 40, maxWidth: 40,
        }
      : {};

  function doDelayedValueUpdate() {
    props.onChange(debounceValue);
  }

  function erase() {
    setSearchValue('');
    props.onChange('');
  }

  return (
    <Container
      style={{ position: 'relative' }}
      border={props.border}
      borderActive={props.borderActive}
      borderHover={props.borderHover}
      TextColor={props.TextColor}
      TextColorHover={props.TextColorHover}
      iconColorHover={props.iconColorHover}
      iconColor={props.iconColor}
    >
      <input
        className="SearchText"
        {...props}
        onFocus={() => setExpand(true)}
        onBlur={() => {
          if (!ui.searchText) setExpand(false);
        }}
        value={searchValue}
        onChange={(e) => {
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
            color: '#B0B7BC',
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
          className="SearchIcon"
          icon="search"
          style={{
            position: 'absolute',
            top: 9,
            right: 9,
            fontSize: 22,
            userSelect: 'none',
            pointerEvents: 'none',
            ...props.iconStyle
          }}
        />
      )}
    </Container>
  );
}

let debounceValue = '';
let inDebounce;
function debounce(func, delay) {
  clearTimeout(inDebounce);
  inDebounce = setTimeout(() => {
    func();
  }, delay);
}

interface inputProps {
  border?: string;
  borderHover?: string;
  borderActive?: string;
  TextColor?: string;
  TextColorHover?: string;
  iconColor?: string;
  iconColorHover?: string;
}

const Container = styled.div<inputProps>`
  .SearchText {
    background: #f2f3f580;
    border: ${(p) => (p.border ? p.border : '1px solid #000')};
    box-sizing: border-box;
    border-radius: 21px;
    padding-left: 20px;
    padding-right: 30px;
    font-style: normal;
    font-weight: 500;
    font-family: Barlow;
    font-size: 16px;
    line-height: 14px;
    height: 35px;
    transition: all 0.4s;
    &::placeholder {
      color: ${(p) => (p.TextColor ? p.TextColor : '#f2f3f580')};
    }
    &:focus {
      border: ${(p) => (p.borderActive ? p.borderActive : '1px solid #000')};
      outline: none;
      caret-color: #a3c1ff;
      &::placeholder {
        color: ${(p) => (p.TextColorHover ? p.TextColorHover : '#f2f3f580')};
      }
    }
  }

  .SearchIcon {
    color: ${(p) => (p.iconColor ? p.iconColor : '#000')};
  }

  &:hover {
    .SearchIcon {
      color: ${(p) => (p.iconColorHover ? p.iconColorHover : '#000')};
    }
    .SearchText {
      border: ${(p) => (p.borderHover ? p.borderHover : '1px solid #000')};
      &:focus {
        border: ${(p) => (p.borderActive ? p.borderActive : '1px solid #000')};
        outline: none;
        caret-color: #a3c1ff;
      }
      &::placeholder {
        color: ${(p) => (p.TextColorHover ? p.TextColorHover : '#f2f3f580')};
      }
    }
  }
  &:active {
    .SearchText {
      border: ${(p) => (p.borderActive ? p.borderActive : '1px solid #000')};
      outline: none;
      caret-color: #a3c1ff;
    }
  }
`;
