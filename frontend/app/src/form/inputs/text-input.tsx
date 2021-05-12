import React from 'react'
import styled from 'styled-components'
import {EuiFormRow, EuiFieldText} from '@elastic/eui'

export default function TextInput(props) {

  return <EuiFormRow label={props.label}>
            <Text name="first" value={props.initialValues.prev_alias}/>
        </EuiFormRow>
}


const Text = styled(EuiFieldText)`

`