import React, { useState } from 'react'
import Widget from './widget'
import FocusedWidget from './focusedWidget'
import FadeLeft from '../../../animated/fadeLeft'
import styled from "styled-components";
import { EuiButton } from '@elastic/eui'

export default function Widgets(props: any) {
    const [selected, setSelected] = useState(null)
    const [showFocused, setShowFocused] = useState(false)

    return <Wrap>
        <FadeLeft
            isMounted={!selected}
            style={{ maxWidth: 500 }}
            dismountCallback={() => setShowFocused(true)}>
            <InnerWrap>
                {props.extras.map((e, i) => {
                    return <Widget key={i} {...e} setSelected={setSelected} />
                })}
            </InnerWrap>
        </FadeLeft>

        <FadeLeft
            isMounted={showFocused}
            dismountCallback={() => setSelected(null)}>
            <>
                <FocusedWidget {...props} setShowFocused={setShowFocused} item={selected} />
            </>
        </FadeLeft>

    </Wrap>

}

const Wrap = styled.div`
`;

const InnerWrap = styled.div`
    display: flex;
    align-content: center;
    justify-content: space-evenly;
    flex-wrap:wrap;
`;