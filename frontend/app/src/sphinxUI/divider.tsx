import React from 'react'
import styled from 'styled-components'

export default function Divider(props: any) {

    return <D style={{ ...props.style }} />
}

const D = styled.div`
height:1px;
background:#EBEDEF;
width:100%;
`

