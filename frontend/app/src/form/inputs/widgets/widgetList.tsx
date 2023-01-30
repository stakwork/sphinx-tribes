import React from 'react';
import styled from 'styled-components';
import Blog from './listItems/blog';
import Offer from './listItems/offer';
import Wanted from './listItems/wanted';
import { EuiButtonIcon } from '@elastic/eui';

export default function WidgetList(props: any) {
  function renderByType(v, i) {
    function wrap(child) {
      return (
        <IWrap
          style={{ cursor: 'pointer' }}
          key={`${i  }listItem`}
          onClick={() => props.setSelected(v, i)}
        >
          {child}
          <Eraser>
            <EuiButtonIcon
              onClick={(e) => {
                e.stopPropagation();
                props.deleteItem(v, i);
              }}
              iconType="trash"
              aria-label="delete"
            />
          </Eraser>
        </IWrap>
      );
    }

    switch (props.schema.class) {
      case 'blog':
        return wrap(<Blog {...v} />);
      case 'offer':
        return wrap(<Offer {...v} />);
      case 'wanted':
        return wrap(<Wanted {...v} />);
      default:
        return <></>;
    }
  }

  return (
    <Wrap>
      <List>
        {props.values &&
          props.values.map((v, i) => {
            return renderByType(v, i);
          })}
      </List>

      {(!props.values || props.values.length < 1) && (
        <IWrap style={{ background: 'none' }}>List is empty</IWrap>
      )}
    </Wrap>
  );
}

export interface IconProps {
  source: string;
}

const Wrap = styled.div`
  color: #fff;
  width: 100%;
`;

const List = styled.div`
  color: #fff;
  width: 100%;
  margin-bottom: 10px;
  display: flex;
  flex-direction: column-reverse;
  align-content: center;
  justify-content: space-evenly;
`;

const IWrap = styled.div`
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
  // border-bottom:1px dashed #1BA9F5;
  padding-bottom: 5px;
  margin: 5px 0;

  background: /* gradient can be an image */ linear-gradient(
      to right,

      #1ba9f5 0%,
      #1ba9f5 100%
    )
    left bottom no-repeat;
  background-size: 100% 1px; /* if linear-gradient, we need to resize it */
`;
// 1BA9F5
// 1d1e24
const Eraser = styled.div`
  cursor: pointer;
`;

const Icon = styled.img<IconProps>`
  background-image: ${(p) => `url(${p.source})`};
  width: 100px;
  height: 100px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: contain; /* Resize the background image to cover the entire container */
`;
