import React from 'react';
import styled from 'styled-components';
import tags from './tags';

const Wrap = styled.div`
  display: flex;
  align-items: center;
  margin-right: 9px;
`;
const Name = styled.span`
  font-size: 10px;
  margin-left: 3px;
`;
const IconWrap = styled.div`
  width: 12px;
  height: 12px;
  border-width: 1px;
  border-style: solid;
  border-radius: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
`;

export default function T(props: { type: string; iconOnly?: boolean }) {
  const { type } = props;
  if (!tags[type]) return <></>;
  const Icon = tags[type].icon;
  const { color } = tags[type];
  return (
    <Wrap className="tag-wrapper">
      <IconWrap style={{ borderColor: color, background: `${color}22` }}>
        <Icon height="8" width="8" />
      </IconWrap>
      {!props.iconOnly && <Name style={{ color: '#556171' }}>{type}</Name>}
    </Wrap>
  );
}
