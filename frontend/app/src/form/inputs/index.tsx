import React from 'react'
import styled from 'styled-components'
import TextInput from './text-input'
import SearchTextInput from './search-text-input'
import TextAreaInput from './text-area-input'
import ImageInput from './img-input'
import GalleryInput from './gallery-input'
import NumberInput from './number-input'
import Widgets from './widgets/index'
import SwitchInput from './switch-input'
import { EuiFormRow, EuiTextArea, EuiFieldText } from '@elastic/eui'
import SelectInput from './select-input'

export default function Input(props: any) {

    function getInput() {
        switch (props.type) {
            case 'text':
                return <TextInput {...props} />
            case 'textarea':
                return <TextAreaInput {...props} />
            case 'search':
                return <SearchTextInput {...props} />
            case 'img':
                return <ImageInput {...props} />
            case 'gallery':
                return <GalleryInput {...props} />
            case 'number':
                return <NumberInput {...props} />
            case 'switch':
                return <SwitchInput {...props} />
            case 'select':
                return <SelectInput {...props} />
            case 'widgets':
                return <Widgets {...props} />
            case 'hidden':
                return <></>
            default:
                return <></>
        }
    }

    return <FieldWrap>
        {getInput()}
    </FieldWrap>
}

const FieldWrap = styled.div`
color:#000 !important;
background-color:#fff !important;
background:#fff !important;
`

export const FieldEnv = styled(EuiFormRow)`
border: 1px solid #DDE1E5;
box-sizing: border-box;
border-radius: 4px;
box-shadow:none !important;
max-width:900px;

.euiFormRow__labelWrapper{
    margin-bottom:0px;
    margin-top:-9px;
    padding-left:10px;
    height:14px;
    label {
        color: #B0B7BC !important;
        background:#ffffff;
    }
}

`
export const FieldText = styled(EuiFieldText)`
background-color:#fff !important;
max-width:900px;
background:#fff !important;
color:${p => p.readOnly ? '#888' : '#000'} !important;
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
    color:#000;
    display:none !important;
}

.euiFormControlLayout--group {
    background-color: #ffffff00 !important;
    transition: none !important;
    display: none !important;
    max-width:900px;
    
}
`
export const FieldTextArea = styled(EuiTextArea)`
background-color:#fff !important;
background:#fff !important;
max-width:900px;
color:#000 !important;
box-shadow:none !important;
`