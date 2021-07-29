import React, { useState } from 'react'
import Widget from './widget'
import FocusedWidget from './focusedWidget'
import FadeLeft from '../../../animated/fadeLeft'
import styled from "styled-components";

export default function Widgets(props: any) {
    const [selected, setSelected] = useState(null)
    const [showFocused, setShowFocused] = useState(false)

    return <Wrap>
        <FadeLeft
            isMounted={!selected}
            style={{ maxWidth: 500 }}
            dismountCallback={() => setShowFocused(true)}>
            <Center>
                <InnerWrap>
                    {props.extras.map((e, i) => {
                        return <Widget parentName={props.name}
                            setFieldValue={props.setFieldValue}
                            values={props.values}
                            key={i} {...e}
                            setSelected={setSelected} />
                    })}
                </InnerWrap>
            </Center>
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
const Center = styled.div`
    display: flex;
    flex:1;
    align-content: center;
    justify-content: center;
`;
const InnerWrap = styled.div`
    display: flex;
    align-content: center;
    justify-content: flex-start;
    flex-wrap:wrap;
    max-width:310px;
`;