import React, { useState, useEffect } from 'react';
import styled from 'styled-components'

const sleep = ms => new Promise(resolve => setTimeout(resolve, ms))

export default function FadeLeft(props) {
    const { drift, isMounted, dismountCallback,
        style, children, alwaysRender, noFadeOnInit,
        direction } = props
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
        if (noFadeOnInit) {
            setOpacity(1)
            setTranslation(0)
        }
    }, [])

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

    if (!alwaysRender && !shouldRender) return <div style={{ position: 'absolute', left: 0, top: 0 }} />
    return (
        <Fader style={{ height: 'inherit', ...style, transform: direction === 'up' ? `translateY(${translation}px)` : `translateX(${translation}px)`, opacity }}>
            {children}
        </Fader>
    );
}

const Fader = styled.div`
  transition:all 0.2s;
`