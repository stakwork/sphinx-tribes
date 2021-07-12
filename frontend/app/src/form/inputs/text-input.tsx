import React from 'react'
import styled from 'styled-components'
import {EuiFormRow, EuiFieldText} from '@elastic/eui'
import type {Props} from './propsType'

export default function TextInput({label, value, handleChange, handleBlur, handleFocus, readOnly, prepend, extraText}:Props) {
  return <>
    <EuiFormRow label={label}>
      <Text name="first" value={value||''} 
        readOnly={readOnly||false}
        onChange={e => handleChange(e.target.value)}
        onBlur={handleBlur}
        onFocus={handleFocus}
        prepend={prepend}
      />
    </EuiFormRow>
    {extraText && value && <ExtraText>{extraText}</ExtraText>}
  </>
}

const Text = styled(EuiFieldText)`

`
const ExtraText = styled.div`
  color:#ddd;
  padding:10px 10px 25px 10px;
  max-width:calc(100% - 20px);
  word-break: break-all;
  font-size:14px;
`