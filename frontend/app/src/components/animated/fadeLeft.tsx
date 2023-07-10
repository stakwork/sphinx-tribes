import { FadeLeftProps } from 'components/interfaces';
import React, { useState, useEffect, useCallback } from 'react';
import styled from 'styled-components';

const sleep = (ms: any) => new Promise((resolve: any) => setTimeout(resolve, ms));

const Fader = styled.div`
  transition: all 0.2s;
`;
const Overlay = styled.div`
  transition: all 0.2s;
  width: 100%;
  height: 100%;
`;
export default function FadeLeft(props: FadeLeftProps) {
  const {
    drift,
    isMounted,
    dismountCallback,
    style,
    children,
    alwaysRender,
    noFadeOnInit,
    direction,
    withOverlay,
    overlayClick,
    noFade
  } = props;
  const [translation, setTranslation] = useState(drift ? drift : -40);
  const [opacity, setOpacity] = useState(0);
  const [shouldRender, setShouldRender] = useState(false);

  function open() {
    setTranslation(0);
    setOpacity(1);
  }

  const close = useCallback(() => {
    setTranslation(drift ? drift : -40);
    setOpacity(0);
  }, [drift]);

  const doAnimation = useCallback(
    (value: any) => {
      if (value === 1) {
        open();
      } else {
        close();
      }
    },
    [close]
  );

  useEffect(() => {
    if (noFadeOnInit) {
      setOpacity(1);
      setTranslation(0);
    }
  }, [noFadeOnInit]);

  useEffect(() => {
    (async () => {
      if (!isMounted) {
        const speed = props.speed ?? 200;
        doAnimation(0);

        await sleep(speed);
        setShouldRender(false);

        if (dismountCallback) dismountCallback();
      } else {
        setShouldRender(true);

        await sleep(5);

        doAnimation(1);
      }
    })();
  }, [dismountCallback, doAnimation, isMounted, props.speed]);

  if (!alwaysRender && !shouldRender)
    return <div style={{ position: 'absolute', left: 0, top: 0 }} />;

  const transformValue =
    direction === 'up' ? `translateY(${translation}px)` : `translateX(${translation}px)`;

  if (withOverlay) {
    return (
      <Overlay style={{ ...style, background: '#000000cc', opacity: !noFade ? opacity : 1 }}>
        <Fader
          id="sphinx-top-level-overlay"
          onClick={(e: any) => {
            let elString: any = e.target;
            elString = elString.outerHTML;
            // close if clicking the overlay
            if (elString && elString.includes('id="sphinx-top-level-overlay"')) {
              if (overlayClick) overlayClick();
              else if (props.close) props.close();
            }
          }}
          style={{ height: 'inherit', ...style, transform: transformValue }}
        >
          {children}
        </Fader>
      </Overlay>
    );
  }
  return (
    <Fader
      style={{
        height: 'inherit',
        ...style,
        transform: transformValue,
        opacity: !noFade ? opacity : 1
      }}
    >
      {children}
    </Fader>
  );
}
