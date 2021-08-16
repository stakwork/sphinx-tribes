import React from 'react'
import styled from 'styled-components'

export default function SearchTextInput(props: any) {
    return <Text
        {...props}
        onChange={e => props.onChange(e.target.value)}
        placeholder={'Search'}
    />
}

const Text = styled.input`
background:#F2F3F580;
border: 1px solid #E0E0E0;
box-sizing: border-box;
border-radius: 21px;
padding-left:20px;
padding-right:10px;
width:100%;
font-style: normal;
font-weight: normal;
font-size: 12px;
line-height: 14px;
height:35px;
`


