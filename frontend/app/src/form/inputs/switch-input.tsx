import React, { useEffect } from 'react'
import styled from 'styled-components'
import { EuiFormRow, EuiSwitch } from '@elastic/eui'
import type { Props } from './propsType'

export default function SwitchInput({ label, value, name, handleChange, handleBlur, handleFocus, readOnly, prepend, extraHTML }: Props) {

    useEffect(() => {
        // if value not initiated, default value true
        if (name === 'show' && value === undefined) handleChange(true)
    }, [])

    return <>
        <EuiFormRow label={label}>
            <EuiSwitch
                label=""
                checked={value}
                onChange={e => {
                    handleChange(e.target.checked)
                }}
                onBlur={handleBlur}
                onFocus={handleFocus}
                compressed
            />
        </EuiFormRow>
        <ExtraText
            style={{ display: value && extraHTML ? 'block' : 'none' }}
            dangerouslySetInnerHTML={{ __html: extraHTML || '' }}
        />
    </>
}

const ExtraText = styled.div`
  color:#ddd;
  padding:10px 10px 25px 10px;
  max-width:calc(100% - 20px);
  word-break: break-all;
  font-size:14px;
`