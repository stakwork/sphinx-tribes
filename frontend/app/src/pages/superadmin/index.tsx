/**
 * Commented out all superadmin restrictions for now
 * To enable colaborations
 */

import React, { useState } from 'react';
import styled from 'styled-components';
// import { useStores } from 'store';
import { MyTable } from './tableComponent';
import { bounties } from './tableComponent/mockBountyData';
import { Header } from './header';
import { Statistics } from './statistics';
import AdminAccessDenied from './accessDenied';

const Container = styled.body`
  height: 100vh; /* Set a fixed height for the container */
  overflow-y: auto; /* Enable vertical scrolling */
  align-items: center;
  margin: 0px auto;
`;

export const SuperAdmin = () => {
  // const { main, ui } = useStores();
  const [isSuperAdmin] = useState(true);

  // const getIsSuperAdmin = useCallback(async () => {
  //   const admin = await main.getSuperAdmin();
  //   setIsSuperAdmin(admin);
  // }, [main]);

  // useEffect(() => {
  //   if (ui.meInfo?.tribe_jwt) {
  //     getIsSuperAdmin();
  //   }
  // }, [main, ui, getIsSuperAdmin]);

  return (
    <>
      {!isSuperAdmin ? (
        <AdminAccessDenied />
      ) : (
        <Container>
          <Header />
          <Statistics />
          <MyTable bounties={bounties} />
        </Container>
      )}
    </>
  );
};
