import React from 'react'
import styled from 'styled-components'
import { EuiIcon } from '@elastic/eui'
import type { Props } from './propsType'
import { FieldEnv, FieldText } from './index'

export default function NumberInput({ error, label, value, handleChange, handleBlur, handleFocus }: Props) {

  let labeltext = label
  if (error) labeltext = labeltext + ` (${error})`

  return <FieldEnv label={labeltext}>
    <R>
      <FieldText name="first" value={value} type="number"
        onChange={e => {
          // dont allow zero or negative numbers
          if (parseInt(e.target.value) < 1) return
          handleChange(e.target.value)
        }}
        onBlur={(e) => {
          // enter 0 on blur if no value
          console.log('onBlur', value)
          if (value === '') handleChange(0)
          handleBlur(e)
        }}
        onFocus={(e) => {
          // remove 0 on focus
          console.log('onFocus', value)
          if (value === 0) handleChange('')
          handleFocus(e)
        }}
      />
      {error && <E>
        <EuiIcon type="alert" size='m' style={{ width: 20, height: 20 }} />
      </E>}
    </R>
  </FieldEnv>
}

const E = styled.div`
  position:absolute;
  right:10px;
  top:0px;
  display:flex;
  height:100%;
  justify-content:center;
  align-items:center;
  color:#45b9f6;
  pointer-events:none;
  user-select:none;
`
const R = styled.div`
  position:relative
`