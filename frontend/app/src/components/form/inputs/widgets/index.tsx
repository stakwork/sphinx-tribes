import React, { useState } from 'react';
import Widget from './widget';
import FocusedWidget from './focusedWidget';
import FadeLeft from '../../../animated/fadeLeft';
import styled from 'styled-components';
import { useEffect } from 'react';

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export default function Widgets(props: any) {
  const [selected, setSelected] = useState(null);
  const [showFocused, setShowFocused] = useState(false);

  useEffect(() => {
    doDelayedScrollTop();
  }, [selected, showFocused]);

  async function doDelayedScrollTop() {
    // we do this so there is no jumping with the animation
    await sleep(140);
    if (props.scrollToTop) props.scrollToTop();
  }
  return (
    <Wrap>
      <FadeLeft
        isMounted={!selected}
        style={{ maxWidth: 500 }}
        dismountCallback={() => setShowFocused(true)}
      >
        <Center>
          <InnerWrap>
            {props.extras.map((e, i) => {
              return (
                <Widget
                  parentName={props.name}
                  setFieldValue={props.setFieldValue}
                  values={props.values}
                  key={i}
                  {...e}
                  setSelected={(e) => {
                    props.setDisableFormButtons(true);
                    setSelected(e);
                  }}
                />
              );
            })}
          </InnerWrap>
        </Center>
      </FadeLeft>

      <FadeLeft isMounted={showFocused} dismountCallback={() => setSelected(null)}>
        <>
          <FocusedWidget {...props} setShowFocused={setShowFocused} item={selected} />
        </>
      </FadeLeft>
    </Wrap>
  );
}

const Wrap = styled.div``;
const Center = styled.div`
  display: flex;
  flex: 1;
  align-content: center;
  justify-content: center;
`;
const InnerWrap = styled.div`
  display: flex;
  align-content: center;
  justify-content: flex-start;
  flex-wrap: wrap;
  max-width: 310px;
`;
