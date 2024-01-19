import React from 'react';
import styled from 'styled-components';
import { EuiFormRow, EuiTextArea, EuiFieldText } from '@elastic/eui';
import { colors } from '../../../config/colors';
import TextInput from './TextInput';
import SearchTextInput from './SearchTextInput';
import TextAreaInput from './TextAreaInput';
import ImageInput from './ImgInput';
import GalleryInput from './GalleryInput';
import NumberInput from './NumberInput';
import Widgets from './widgets/index';
import SwitchInput from './SwitchInput';
import LoomVideoInput from './LoomVideoInput';
import SelectInput from './SelectInput';
import SearchableSelectInput from './SearchableSelectInput';
import MultiSelectInput from './MultiSelectInput';
import CreatableMultiSelectInput from './CreatableMultiSelectInput';
import Date from './widgets/Date';
import InvitePeopleSearch from './widgets/PeopleSearch';
import LoomVideoInputNew from './LoomVideoInputNew';
import TextInputNew from './TextInputNew2';
import NumberInputNew from './NumberInputNew';
import NumberSatsInput from './NumberSatsInput';
import TextAreaInputNew from './TextAreaInputNew';
import CreatableMultiSelectInputNew from './CreatableMultiSelectInputNew';

interface styledProps {
  color?: any;
  height?: any;
  width?: any;
  isTextField?: any;
}

const FieldWrap = styled.div<styledProps>`
  color: ${(p: any) => p?.color && p.color.pureBlack} !important;
  background-color: ${(p: any) => p?.color && p.color.pureWhite} !important;
  background: ${(p: any) => p?.color && p.color.pureWhite} !important;
`;

export const Note = styled.div<styledProps>`
  display: flex;
  padding-left: 10px;
  flex-wrap: wrap;
  font-size: 12px;
  color: ${(p: any) => p?.color && p.color.grayish.G60A};
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
  border: ${(p: any) =>
    p.border === 'bottom'
      ? ''
      : p?.isTop
      ? `1px solid ${p?.color && p.color.pureWhite}`
      : `1px solid ${p?.color && p.color.grayish.G600}`};
  border-bottom: ${(p: any) => (p.border === 'bottom' ? `1px solid ${p.color.grayish.G600}` : '')};
  box-sizing: border-box;
  border-radius: ${(p: any) => (p.border === 'bottom' ? '0px' : '4px')};
  box-shadow: none !important;
  max-width: 900px;
  min-height: ${(p: any) => (p?.isTextField ? '40' : '')};
  max-height: ${(p: any) => (p?.isTextField ? '40' : '')};
  margin-bottom: 24px;
  border: ${(p: any) => p?.isFill && `1px solid ${p?.color && p.color.grayish.G600}`};
  .euiFormRow__labelWrapper {
    margin-bottom: -20px;
    margin-top: ${(p: any) => (p?.isTextField ? '12px' : '10px')};
    padding-left: 10.5px;
    height: ${(p: any) => (p?.isTextField ? '6px' : '14px')};
    position: relative;
    label {
      color: ${(p: any) => p?.color && p.color.grayish.G300} !important;
      background: ${(p: any) => p?.color && p.color.pureWhite};
      z-index: 10;
      position: ${(p: any) => p?.isTop && 'absolute'};
      top: ${(p: any) => p?.isTop && '-20px'} !important;
      font-family: 'Barlow';
      font-style: normal;
      font-weight: 500;
      font-size: 14px;
      display: flex;
      align-items: left;
    }
  }
`;
export const FieldText = styled(EuiFieldText)<styledProps>`
  background-color: ${(p: any) => p?.color && p.color.pureWhite} !important;
  max-width: 900px;
  background: ${(p: any) => p?.color && p.color.pureWhite} !important;
  color: ${(p: any) =>
    p.readOnly ? `${p.color.grayish.G60A}` : `${p.color.pureBlack}`} !important;
  box-shadow: none !important;
  height: ${(p: any) => (p?.isTextField ? '12px' : '')};
  margin-top: ${(p: any) => (p?.isTextField ? '8px' : '')};

  .euiFormRow__labelWrapper .euiFormControlLayout--group {
    background-color: ${(p: any) => p?.color && p.color.pureWhite} !important;
    background: ${(p: any) => p?.color && p.color.pureWhite} !important;
    box-shadow: none !important;
  }

  .euiFormRow__fieldWrapper .euiFormControlLayout {
    background-color: ${(p: any) => p?.color && p.color.pureWhite} !important;
    background: ${(p: any) => p?.color && p.color.pureWhite} !important;
    box-shadow: none !important;
  }

  .euiFormLabel euiFormControlLayout__prepend {
    background-color: ${(p: any) => p?.color && p.color.pureWhite} !important;
    background: ${(p: any) => p?.color && p.color.pureWhite} !important;
    box-shadow: none !important;
    color: ${(p: any) => p?.color && p.color.grayish.G71} !important;
    display: none !important;
    line-height: 1rem;
  }

  .euiFormControlLayout--group {
    background-color: #ffffff00 !important;
    transition: none !important;
    display: none !important;
    max-width: 900px;
  }
`;
export const FieldTextArea = styled(EuiTextArea)<styledProps>`
  width: ${(p: any) => p?.color && p.color.width};
  background-color: ${(p: any) => p?.color && p.color.pureWhite} !important;
  background: ${(p: any) => p?.color && p.color.pureWhite} !important;
  max-width: 900px;
  width: 292px;
  min-height: 167px;
  color: ${(p: any) => p?.color && p.color.pureBlack} !important;
  box-shadow: none !important;
  line-height: 17.6px;
`;
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
      case 'numbersats':
        return <NumberSatsInput {...props} />;
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
