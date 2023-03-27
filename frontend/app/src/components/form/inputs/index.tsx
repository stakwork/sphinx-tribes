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
import { colors } from '../../../config/colors';
import LoomVideoInputNew from './loom-video-input-new';
import TextInputNew from './text-input-new';
import NumberInputNew from './number-input-new';
import TextAreaInputNew from './text-area-input-new';
import CreatableMultiSelectInputNew from './creatable-multi-select-input-new';

export default function Input(props: any) {
  const color = colors['light'];
  function getInput() {
    switch (props.type) {
      case 'space':
        return <div style={{ height: 10 }} />;
      case 'text':
        return props?.newDesign ? <TextInputNew {...props} /> : <TextInput {...props} />;
      case 'textarea':
        return props?.newDesign ? <TextAreaInputNew {...props} /> : <TextAreaInput {...props} />;
      case 'search':
        return <SearchTextInput {...props} />;
      case 'img':
        return <ImageInput {...props} />;
      case 'imgcanvas':
        return <ImageInput notProfilePic={true} {...props} />;
      case 'gallery':
        return <GalleryInput {...props} />;
      case 'number':
        return props?.newDesign ? <NumberInputNew {...props} /> : <NumberInput {...props} />;
      case 'loom':
        return props?.newDesign ? <LoomVideoInputNew {...props} /> : <LoomVideoInput {...props} />;
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
        return props?.newDesign ? (
          <CreatableMultiSelectInputNew {...props} />
        ) : (
          <CreatableMultiSelectInput {...props} />
        );
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
  height?: any;
  width?: any;
  isTextField?: any;
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
  height?: any;
  width?: any;
  isTextField?: any;
}

export const FieldEnv = styled(EuiFormRow as any)<fieldEnvProps>`
  border: ${(p) =>
    p.border === 'bottom'
      ? ''
      : p?.isTop
      ? `1px solid ${p?.color && p.color.pureWhite}`
      : `1px solid ${p?.color && p.color.grayish.G600}`};
  border-bottom: ${(p) => (p.border === 'bottom' ? `1px solid ${p.color.grayish.G600}` : '')};
  box-sizing: border-box;
  border-radius: ${(p) => (p.border === 'bottom' ? '0px' : '4px')};
  box-shadow: none !important;
  max-width: 900px;
  min-height: ${(p) => (p?.isTextField ? '40' : '')};
  max-height: ${(p) => (p?.isTextField ? '40' : '')};
  margin-bottom: 24px;
  border: ${(p) => p?.isFill && `1px solid ${p?.color && p.color.grayish.G600}`};
  .euiFormRow__labelWrapper {
    margin-bottom: -20px;
    margin-top: ${(p) => (p?.isTextField ? '12px' : '10px')};
    padding-left: 16px;
    height: ${(p) => (p?.isTextField ? '6px' : '14px')};
    position: relative;
    label {
      color: ${(p) => p?.color && p.color.grayish.G300} !important;
      background: ${(p) => p?.color && p.color.pureWhite};
      z-index: 10;
      position: ${(p) => p?.isTop && 'absolute'};
      top: ${(p) => p?.isTop && '-20px'} !important;
      font-family: Barlow;
      font-style: normal;
      font-weight: 500;
      font-size: 14px;
      display: flex;
      align-items: center;
    }
  }
`;
export const FieldText = styled(EuiFieldText)<styledProps>`
background-color: ${(p) => p?.color && p.color.pureWhite} !important;
max-width:900px;
background:${(p) => p?.color && p.color.pureWhite} !important;
color:${(p) => (p.readOnly ? `${p.color.grayish.G60A}` : `${p.color.pureBlack}`)} !important;
box-shadow: none !important;
height: ${(p) => (p?.isTextField ? '12px' : '')}  ;
margin-top: ${(p) => (p?.isTextField ? '2px' : '')}; 

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
    line-height: 1rem;
}

.euiFormControlLayout--group {
    background-color: #ffffff00 !important;
    transition: none !important;
    display: none !important;
    max-width:900px;
    
}


`;
export const FieldTextArea = styled(EuiTextArea)<styledProps>`
  // min-height: ${(p) => p?.height && p.height} !important;
  // max-height: ${(p) => p?.width && p.height} !important;
  width: ${(p) => p?.color && p.color.width};
  background-color: ${(p) => p?.color && p.color.pureWhite} !important;
  background: ${(p) => p?.color && p.color.pureWhite} !important;
  max-width: 900px;
  color: ${(p) => p?.color && p.color.pureBlack} !important;
  box-shadow: none !important;
  // border-bottom: ${(p) => p?.color && `1px solid ${p.color.grayish.G600}`};
  line-height: 17.6px;
`;
