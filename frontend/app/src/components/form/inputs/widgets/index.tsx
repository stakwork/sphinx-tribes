import React, { useState } from 'react';
import styled from 'styled-components';
import { useEffect } from 'react';
import FadeLeft from '../../../animated/FadeLeft';
import Widget from './Widget';
import FocusedWidget from './FocusedWidget';

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
async function sleep(ms: number) {
  return new Promise((resolve: any) => setTimeout(resolve, ms));
}

export default function Widgets(props: any) {
  const [selected, setSelected] = useState(null);
  const [showFocused, setShowFocused] = useState(false);

  async function doDelayedScrollTop() {
    // we do this so there is no jumping with the animation
    await sleep(140);
    if (props.scrollToTop) props.scrollToTop();
  }
  useEffect(() => {
    doDelayedScrollTop();
  }, [selected, showFocused]);

  return (
    <Wrap>
      <FadeLeft
        isMounted={!selected}
        style={{ maxWidth: 500 }}
        dismountCallback={() => setShowFocused(true)}
        drift={0}
        withOverlay={false}
      >
        <Center>
          <InnerWrap>
            {props.extras.map((e: any, i: number) => (
              <Widget
                parentName={props.name}
                setFieldValue={props.setFieldValue}
                values={props.values}
                key={i}
                {...e}
                setSelected={(e: any) => {
                  props.setDisableFormButtons(true);
                  setSelected(e);
                }}
              />
            ))}
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
