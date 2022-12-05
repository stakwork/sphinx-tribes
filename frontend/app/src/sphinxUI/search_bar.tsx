import MaterialIcon from '@material/react-material-icon';
import React, { useState } from 'react';
import styled from 'styled-components';
import { colors } from '../colors';
import { useStores } from '../store';

export default function SearchTextInput(props: any) {
  const color = colors['light'];
  const { ui } = useStores();
  const [searchValue, setSearchValue] = useState(ui.searchText || '');
  const [expand, setExpand] = useState(ui.searchText ? true : false);

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
      color={color}>
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
        placeholder={props.placeholder}
        style={{ ...props.style }}
      />
      {searchValue ? (
        <MaterialIcon
          icon="close"
          onClick={() => erase()}
          style={{
            position: 'absolute',
            color: color.grayish.G300,
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
  color?: any;
}

const Container = styled.div<inputProps>`
  .SearchText {
    background: ${(p) => p.color && p.color.grayish.G600} !important;
    border: ${(p) => (p.border ? p.border : `1px solid ${p.color.pureBlack}`)};
    box-sizing: border-box;
    border-radius: 200px;
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
      color: ${(p) => (p.TextColor ? p.TextColor : `${p.color.grayish.G65}`)};
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 400;
      font-size: 16px;
      line-height: 19px;
    }
    &:focus {
      border: ${(p) => (p.borderActive ? p.borderActive : `1px solid ${p.color.pureBlack}`)};
      outline: none;
      caret-color: ${(p) => p.color && p.color.light_blue100};
      &::placeholder {
        color: ${(p) => (p.TextColorHover ? p.TextColorHover : `${p.color.grayish.G65}`)};
      }
    }
    &:focus-within {
      background: ${(p) => p.color && p.color.grayish.G950} !important;
    }
  }

  .SearchIcon {
    color: ${(p) => (p.iconColor ? p.iconColor : `${p.color.pureBlack}`)};
  }

  &:hover {
    .SearchIcon {
      color: ${(p) => (p.iconColorHover ? p.iconColorHover : `${p.color.pureBlack}`)};
    }
    .SearchText {
      border: ${(p) => (p.borderHover ? p.borderHover : `1px solid ${p.color.pureBlack}`)};
      &:focus {
        border: ${(p) => (p.borderActive ? p.borderActive : `1px solid ${p.color.pureBlack}`)};
        outline: none;
        caret-color: ${(p) => p.color && p.color.light_blue100};
      }
      &::placeholder {
        color: ${(p) => (p.TextColorHover ? p.TextColorHover : `${p.color.grayish.G65}`)};
      }
    }
  }
  &:active {
    .SearchText {
      border: ${(p) => (p.borderActive ? p.borderActive : `1px solid ${p.color.pureBlack}`)};
      outline: none;
      caret-color: ${(p) => p.color && p.color.light_blue100};
    }
  }
`;
