import React from 'react'
import styled from 'styled-components'
import { EuiFormRow, EuiTextArea } from '@elastic/eui'
import type { Props } from './propsType'

export default function TextAreaInput({ label, value, handleChange, handleBlur, handleFocus, readOnly, prepend, extraHTML }: Props) {
    // console.log("TEXTAREA", label, extraHTML)
    return <>
        <EuiFormRow label={label}>
            <Text name="first" value={value || ''}
                readOnly={readOnly || false}
                onChange={e => handleChange(e.target.value)}
                onBlur={handleBlur}
                onFocus={handleFocus}
            // prepend={prepend}
            />
        </EuiFormRow>
        <ExtraText
            style={{ display: value && extraHTML ? 'block' : 'none' }}
            dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
        />
    </>
}

const Text = styled(EuiTextArea)`

`
const ExtraText = styled.div`
  color:#ddd;
  padding:10px 10px 25px 10px;
  max-width:calc(100% - 20px);
  word-break: break-all;
  font-size:14px;
`