import React from 'react';
import styled from 'styled-components';
import TextInput from './text-input';
import SearchTextInput from './search-text-input';
import TextAreaInput from './text-area-input';
import ImageInput from './img-input';
import GalleryInput from './gallery-input';
import NumberInput from './number-input';
import Widgets from './widgets/index';
import SwitchInput from './switch-input';
import LoomVideoInput from './loom-video-input';
import { EuiFormRow, EuiTextArea, EuiFieldText } from '@elastic/eui';
import SelectInput from './select-input';
import SearchableSelectInput from './searchable-select-input';
import MultiSelectInput from './multi-select-input';
import CreatableMultiSelectInput from './creatable-multi-select-input';
import Date from './widgets/date';
import InvitePeopleSearch from './widgets/PeopleSearch';
import { colors } from '../../colors';

export default function Input(props: any) {
  const color = colors['light'];
  function getInput() {
    switch (props.type) {
      case 'space':
        return <div style={{ height: 10 }} />;
      case 'text':
        return <TextInput {...props} />;
      case 'textarea':
        return <TextAreaInput {...props} />;
      case 'search':
        return <SearchTextInput {...props} />;
      case 'img':
        return <ImageInput {...props} />;
      case 'imgcanvas':
        return <ImageInput notProfilePic={true} {...props} />;
      case 'gallery':
        return <GalleryInput {...props} />;
      case 'number':
        return <NumberInput {...props} />;
      case 'loom':
        return <LoomVideoInput {...props} />;
      case 'switch':
        return <SwitchInput {...props} />;
      case 'select':
        return <SelectInput {...props} />;
      case 'searchableselect':
        return props?.newDesign ? (
          <InvitePeopleSearch {...props} />
        ) : (
          <SearchableSelectInput {...props} />
        );
      case 'multiselect':
        return <MultiSelectInput {...props} />;
      case 'creatablemultiselect':
        return <CreatableMultiSelectInput {...props} />;
      case 'widgets':
        return <Widgets {...props} />;
      case 'date':
        return <Date {...props} />;
      case 'hidden':
        return <></>;
      default:
        return <></>;
    }
  }

  return <FieldWrap color={color}>{getInput()}</FieldWrap>;
}

interface styledProps {
  color?: any;
}

const FieldWrap = styled.div<styledProps>`
  color: ${(p) => p?.color && p.color.pureBlack} !important;
  background-color: ${(p) => p?.color && p.color.pureWhite} !important;
  background: ${(p) => p?.color && p.color.pureWhite} !important;
`;

export const Note = styled.div<styledProps>`
  display: flex;
  padding-left: 10px;
  flex-wrap: wrap;
  font-size: 12px;
  color: ${(p) => p?.color && p.color.grayish.G60A};
  margin-bottom: 16px;
  margin-top: -6px;
  max-width: 400px;
  font-style: italic;
`;

interface fieldEnvProps {
  readonly border: string;
  isTop?: boolean;
  isFill?: boolean;
  color?: any;
}

export const FieldEnv = styled(EuiFormRow as any)<fieldEnvProps>`
  border: ${(p) =>
    p.border === 'bottom'
      ? ''
      : p?.isTop
      ? `1px solid ${p.color.pureWhite}`
      : `1px solid ${p.color.grayish.G600}`};
  border-bottom: ${(p) => (p.border === 'bottom' ? `1px solid ${p.color.grayish.G600}` : '')};
  box-sizing: border-box;
  border-radius: ${(p) => (p.border === 'bottom' ? '0px' : '4px')};
  box-shadow: none !important;
  max-width: 900px;
  margin-bottom: 24px;
  border: ${(p) => p?.isFill && `1px solid ${p.color.grayish.G600}`};
  .euiFormRow__labelWrapper {
    margin-bottom: -20px;
    margin-top: 10px;
    padding-left: 10px;
    height: 14px;
    position: relative;
    label {
      color: ${(p) => p.color && p.color.grayish.G300} !important;
      background: ${(p) => p?.color && p.color.pureWhite};
      z-index: 10;
      position: ${(p) => p?.isTop && 'absolute'};
      top: ${(p) => p?.isTop && '-20px'};
    }
  }
`;
export const FieldText = styled(EuiFieldText)<styledProps>`
background-color: ${(p) => p?.color && p.color.pureWhite} !important;
max-width:900px;
background:${(p) => p?.color && p.color.pureWhite} !important;
color:${(p) => (p.readOnly ? `${p.color.grayish.G60A}` : `${p.color.pureBlack}`)} !important;
box-shadow:none !important;

.euiFormRow__labelWrapper .euiFormControlLayout--group{
    background-color:${(p) => p?.color && p.color.pureWhite} !important;
    background:${(p) => p?.color && p.color.pureWhite} !important;
    box-shadow:none !important;
}

.euiFormRow__fieldWrapper .euiFormControlLayout {
    background-color:${(p) => p?.color && p.color.pureWhite} !important;
    background:${(p) => p?.color && p.color.pureWhite} !important;
    box-shadow:none !important;
}

.euiFormLabel euiFormControlLayout__prepend{
    background-color:${(p) => p?.color && p.color.pureWhite} !important;
    background:${(p) => p?.color && p.color.pureWhite} !important;
    box-shadow:none !important;
    color: ${(p) => p?.color && p.color.grayish.G71} !important;s
    display:none !important;
}

.euiFormControlLayout--group {
    background-color: #ffffff00 !important;
    transition: none !important;
    display: none !important;
    max-width:900px;
    
}


`;
export const FieldTextArea = styled(EuiTextArea)<styledProps>`
  background-color: ${(p) => p?.color && p.color.pureWhite} !important;
  background: ${(p) => p?.color && p.color.pureWhite} !important;
  max-width: 900px;
  color: ${(p) => p?.color && p.color.pureBlack} !important;
  box-shadow: none !important;
`;
