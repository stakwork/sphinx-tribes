import React, { useState, useEffect } from 'react';
import styled from 'styled-components'

const sleep = ms => new Promise(resolve => setTimeout(resolve, ms))

export default function FadeLeft(props) {
    const { drift, isMounted, dismountCallback, style, children } = props
    const [translation, setTranslation] = useState(drift ? drift : -40)
    const [opacity, setOpacity] = useState(0)
    const [shouldRender, setShouldRender] = useState(false)

    function open() {
        setTranslation(0)
        setOpacity(1)
    }

    function close() {
        setTranslation(drift ? drift : -40)
        setOpacity(0)
    }

    function doAnimation(value) {
        if (value === 1) {
            open()
        } else {
            close()
        }
    }

    useEffect(() => {
        (async () => {
            if (!isMounted) {
                let speed = 200
                if (props.speed) speed = props.speed
                doAnimation(0)
                await sleep(speed)
                console.log('unmount')
                setShouldRender(false)
                if (dismountCallback) dismountCallback()
            } else {
                setShouldRender(true);
                await sleep(5)
                console.log('render true')
                doAnimation(1)
            }
        })()
    }, [isMounted])

    if (!shouldRender) return <div style={{ position: 'absolute', left: 0, top: 0 }} />
    return (
        <Fader style={{ ...style, transform: `translateX(${translation}px)`, opacity }}>
            {children}
        </Fader>
    );
}

const Fader = styled.div`
  transition:all 0.2s;
`