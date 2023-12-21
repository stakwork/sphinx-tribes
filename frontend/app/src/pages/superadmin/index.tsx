import React from 'react';
import styled from 'styled-components';
import { MyTable } from './tableComponent';
import { bounties } from './tableComponent/mockBountyData';
import { Header } from './header';
import { Statistics } from './statistics';

const Container = styled.body`
  height: 100vh; /* Set a fixed height for the container */
  overflow-y: auto; /* Enable vertical scrolling */
  width: 1366px;
  align-items: center;
  margin: 0px auto;
`;

export const SuperAdmin = () => (
  <Container>
    <Header />
    <Statistics />
    <MyTable bounties={bounties} />
  </Container>
);
