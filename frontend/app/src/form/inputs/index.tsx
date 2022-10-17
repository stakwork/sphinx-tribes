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

export default function Input(props: any) {
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
        return <SearchableSelectInput {...props} />;
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

  return <FieldWrap>{getInput()}</FieldWrap>;
}

const FieldWrap = styled.div`
  color: #000 !important;
  background-color: #fff !important;
  background: #fff !important;
`;

export const Note = styled.div`
  display: flex;
  padding-left: 10px;
  flex-wrap: wrap;
  font-size: 12px;
  color: #888;
  margin-bottom: 16px;
  margin-top: -6px;
  max-width: 400px;
  font-style: italic;
`;

interface fieldEnvProps {
  readonly border: string;
}

export const FieldEnv = styled(EuiFormRow as any)<fieldEnvProps>`
  border: ${(p) => (p.border === 'bottom' ? '' : '1px solid #dde1e5')};
  border-bottom: ${(p) => (p.border === 'bottom' ? '1px solid #dde1e5' : '')};
  box-sizing: border-box;
  border-radius: ${(p) => (p.border === 'bottom' ? '0px' : '4px')};
  box-shadow: none !important;
  max-width: 900px;

  .euiFormRow__labelWrapper {
    margin-bottom: 0px;
    margin-top: -9px;
    padding-left: 10px;
    height: 14px;
    label {
      color: #b0b7bc !important;
      background: #ffffff;
      z-index: 10;
    }
  }
`;
export const FieldText = styled(EuiFieldText)`
background-color:#fff !important;
max-width:900px;
background:#fff !important;
color:${(p) => (p.readOnly ? '#888' : '#000')} !important;
box-shadow:none !important;

.euiFormRow__labelWrapper .euiFormControlLayout--group{
    background-color:#fff !important;
    background:#fff !important;
    box-shadow:none !important;
}

.euiFormRow__fieldWrapper .euiFormControlLayout {
    background-color:#fff !important;
    background:#fff !important;
    box-shadow:none !important;
}

.euiFormLabel euiFormControlLayout__prepend{
    background-color:#fff !important;
    background:#fff !important;
    box-shadow:none !important;
    color:#999 !important;s
    display:none !important;
}

.euiFormControlLayout--group {
    background-color: #ffffff00 !important;
    transition: none !important;
    display: none !important;
    max-width:900px;
    
}


`;
export const FieldTextArea = styled(EuiTextArea)`
  background-color: #fff !important;
  background: #fff !important;
  max-width: 900px;
  color: #000 !important;
  box-shadow: none !important;
`;
