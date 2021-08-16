import React, { useEffect, useState } from 'react'
import styled from 'styled-components'

import { colors } from '../../colors'
import Checkbox from '../../form/inputs/checkbox-input'

export default function TagComponent(props: any) {

    const c = colors['light']

    return <Tag>

        <Env>
            <Checkbox {...props} />
            <Name>name</Name>
        </Env>
        <Number>234</Number>
    </Tag>

}

const Tag = styled.div`
            display:flex;
            align-items:center;
            width:100%;
            justify-content:space-between;
            margin-bottom:5px;
            `

const Env = styled.div`
            display:flex;
            align-items:center;
            `

const Name = styled.div`
            font-style: normal;
            font-weight: normal;
            font-size: 12px;
            line-height: 31px;
            margin-left:10px;

            `

const Number = styled.div`
            font-style: normal;
            font-weight: normal;
            font-size: 12px;
            line-height: 31px;
            text-align: right;

            color: #B0B7BC;
            `