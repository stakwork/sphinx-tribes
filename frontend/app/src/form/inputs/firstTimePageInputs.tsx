import React from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiTextArea, EuiFieldText } from '@elastic/eui';
import SelectInput from './select-input';
import TextInput from '../inputs_new/text-input';
import TextAreaInput from '../inputs_new/text-area-input';
import SearchTextInput from '../inputs_new/search-text-input';
import ImageInput from '../inputs_new/img-input';
import GalleryInput from '../inputs_new/gallery-input';
import NumberInput from '../inputs_new/number-input';
import LoomVideoInput from '../inputs_new/loom-video-input';
import SwitchInput from '../inputs_new/switch-input';
import SearchableSelectInput from '../inputs_new/searchable-select-input';
import MultiSelectInput from '../inputs_new/multi-select-input';
import CreatableMultiSelectInput from '../inputs_new/creatable-multi-select-input';
import Widgets from '../inputs_new/widgets';
import Date from '../inputs_new/widgets/date';

export default function FirstTimeInput(props: any) {
  console.log(props);
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

export const FieldEnv = styled(EuiFormRow as any)`
  border-bottom: 1px solid #dde1e5;
  box-sizing: border-box;
  //   border-radius: 4px;
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
