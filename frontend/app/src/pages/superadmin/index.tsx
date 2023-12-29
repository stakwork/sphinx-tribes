/**
 * Commented out all superadmin restrictions for now
 * To enable colaborations
 */

import React, { useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';
import { useStores } from 'store';
import moment from 'moment';
import { MyTable } from './tableComponent';
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
  //Todo: Remove all comments when metrcis development is done
  const { main } = useStores();
  const [isSuperAdmin] = useState(true);
  const [bounties, setBounties] = useState<any[]>([]);

  /**
   * Todo use the same date range,
   * and status for all child components
   * */
  const [endDate] = useState(moment().unix());
  const [startDate] = useState(moment().subtract(30, 'days').unix());

  // const getIsSuperAdmin = useCallback(async () => {
  //   const admin = await main.getSuperAdmin();
  //   setIsSuperAdmin(admin);
  // }, [main]);

  // useEffect(() => {
  //   if (ui.meInfo?.tribe_jwt) {
  //     getIsSuperAdmin();
  //   }
  // }, [main, ui, getIsSuperAdmin]);

  const getBounties = useCallback(async () => {
    if (startDate && endDate) {
      const bounties = await main.getBountiesByRange(String(startDate), String(endDate));
      setBounties(bounties);
    }
  }, [main, startDate, endDate]);

  useEffect(() => {
    getBounties();
  }, [getBounties]);

  return (
    <>
      {!isSuperAdmin ? (
        <AdminAccessDenied />
      ) : (
        <Container>
          <Header />
          <Statistics />
          <MyTable bounties={bounties} startDate={startDate} endDate={endDate} />
        </Container>
      )}
    </>
  );
};
