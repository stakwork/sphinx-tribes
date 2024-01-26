import React from 'react';
import styled from 'styled-components';
import { EuiButtonIcon } from '@elastic/eui';
import Blog from './listItems/Blog';
import Offer from './listItems/Offer';
import Wanted from './listItems/Wanted';
import { WidgetListProps } from './interfaces';

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
export default function WidgetList(props: WidgetListProps) {
  function renderByType(v: any, i: any) {
    function wrap(child: any) {
      return (
        <IWrap
          style={{ cursor: 'pointer' }}
          key={`${i}listItem`}
          onClick={() => props.setSelected(v, i)}
        >
          {child}
          <Eraser>
            <EuiButtonIcon
              onClick={(e: any) => {
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
      <List>{props.values && props.values.map((v: any, i: number) => renderByType(v, i))}</List>

      {(!props.values || props.values.length < 1) && (
        <IWrap style={{ background: 'none' }}>List is empty</IWrap>
      )}
    </Wrap>
  );
}
