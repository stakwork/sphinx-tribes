import React from 'react';
import styled from 'styled-components';
import { BlogPost } from '../interfaces';

const Wrap = styled.div`
  display: flex;
`;
export default function Blog(props: BlogPost) {
  return (
    <Wrap>
      <div>{props.title}</div>
      <div>{props.created}</div>
    </Wrap>
  );
}
