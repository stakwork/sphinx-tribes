import React from 'react';
import styled from 'styled-components';
import { BlogPost } from '../interfaces';

export default function Blog(props: BlogPost) {
  return (
    <Wrap>
      <div>{props.title}</div>
      <div>{props.created}</div>
    </Wrap>
  );
}

const Wrap = styled.div`
  display: flex;
`;
