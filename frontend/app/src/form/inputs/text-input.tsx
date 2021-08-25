import React from 'react'
import styled from 'styled-components'
import { EuiFormRow, EuiFieldText, EuiIcon } from '@elastic/eui'
import type { Props } from './propsType'
import { FieldEnv, FieldText } from './index'

export default function TextInput({ error, label, value, handleChange, handleBlur, handleFocus, readOnly, prepend, extraHTML }: Props) {

  let labeltext = label
  if (error) labeltext = labeltext + ` (${error})`

  return <>
    <FieldEnv label={labeltext}>
      <R>
        <FieldText name="first" value={value || ''}
          readOnly={readOnly || false}
          onChange={e => handleChange(e.target.value)}
          onBlur={handleBlur}
          onFocus={handleFocus}
          prepend={prepend}
        />
        {error && <E>
          <EuiIcon type="alert" size='m' style={{ width: 20, height: 20 }} />
        </E>}
      </R>
    </FieldEnv>
    <ExtraText
      style={{ display: value && extraHTML ? 'block' : 'none' }}
      dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
    />
  </>
}

const ExtraText = styled.div`
  padding:2px 10px 25px 10px;
  max-width:calc(100% - 20px);
  word-break: break-all;
  font-size:14px;
`

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

