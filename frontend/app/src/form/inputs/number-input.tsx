import React from 'react'
import styled from 'styled-components'
import {EuiFormRow, EuiFieldText} from '@elastic/eui'
import type {Props} from './propsType'

export default function NumberInput({label, value, handleChange, handleBlur, handleFocus}:Props) {
  return <EuiFormRow label={label}>
    <Text name="first" value={value} type="number"
      onChange={e => handleChange(e.target.value)}
      onBlur={handleBlur}
      onFocus={handleFocus}
    />
  </EuiFormRow>
}


const Text = styled(EuiFieldText)`

`