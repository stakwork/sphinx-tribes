/**
 * Commented out all superadmin restrictions for now
 * To enable colaborations
 */
import React, { useCallback, useEffect, useState } from 'react';
import { EuiLoadingSpinner } from '@elastic/eui';
import styled from 'styled-components';
import { BountyMetrics, BountyStatus } from 'store/main';
import { useStores } from 'store';
import moment from 'moment';
import { useInViewPort } from 'hooks';
import { MyTable } from './tableComponent';
import { Header } from './header';
import { Statistics } from './statistics';
import AdminAccessDenied from './accessDenied';

const Container = styled.body`
  height: 100vh; /* Set a fixed height for the container */
  overflow-y: auto; /* Enable vertical scrolling */
  align-items: center;
  margin: 0px auto;
  padding: 4.5rem 0;
`;

const LoaderContainer = styled.div`
  height: 20%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const SuperAdmin = () => {
  //Todo: Remove all comments when metrcis development is done
  const { main } = useStores();
  const [isSuperAdmin] = useState(true);
  const [bounties, setBounties] = useState<any[]>([]);
  const [bountyMetrics, setBountyMetrics] = useState<BountyMetrics | undefined>(undefined);
  const [bountyStatus, setBountyStatus] = useState<BountyStatus>({
    Open: false,
    Assigned: false,
    Paid: false
  });
  const [dropdownValue, setDropdownValue] = useState('all');
  const [loading, setLoading] = useState(false);

  /**
   * Todo use the same date range,
   * and status for all child components
   * */

  const [endDate, setEndDate] = useState(moment().unix());
  const [startDate, setStartDate] = useState(moment().subtract(30, 'days').unix());

  
  const [inView, ref] = useInViewPort({
    rootMargin: '0px',
    threshold: 0.25
  });

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
    setLoading(true);
    if (startDate && endDate) {
      try {
        const bounties = await main.getBountiesByRange(
          {
            start_date: String(startDate),
            end_date: String(endDate)
          },
          {
            resetPage: true,
            ...bountyStatus
          }
        );
        setBounties(bounties);
      } catch (error) {
        // Handle errors if any
        console.error('Error fetching total bounties:', error);
      } finally {
        // Set loading to false regardless of success or failure
        setLoading(false);
      }
    }
  }, [main, startDate, endDate, bountyStatus]);

  useEffect(() => {
    getBounties();
  }, [getBounties]);

  const normalizeMetrics = (data: any): BountyMetrics => ({
    bounties_posted: data.BountiesPosted || data.bounties_posted,
    bounties_paid: data.BountiesPaid || data.bounties_paid,
    bounties_paid_average: data.bounties_paid_average || data.BountiesPaidPercentage,
    sats_posted: data.sats_posted || data.SatsPosted,
    sats_paid: data.sats_paid || data.SatsPaid,
    sats_paid_percentage: data.sats_paid_percentage || data.SatsPaidPercentage,
    average_paid: data.average_paid || data.AveragePaid,
    average_completed: data.average_completed || data.AverageCompleted
  });

  const getMetrics = useCallback(async () => {
    if (startDate && endDate) {
      try {
        const metrics = await main.getBountyMetrics(String(startDate), String(endDate));
        const normalizedMetrics = normalizeMetrics(metrics);
        setBountyMetrics(normalizedMetrics);
      } catch (error) {
        console.error('Error fetching metrics:', error);
      }
    }
  }, [main, startDate, endDate]);

  useEffect(() => {
    getMetrics();
  }, [getMetrics]);

  return (
    <>
      {!isSuperAdmin ? (
        <AdminAccessDenied />
      ) : (
        <Container>
          <Header
            startDate={startDate}
            endDate={endDate}
            setStartDate={setStartDate}
            setEndDate={setEndDate}
          />
          <Statistics freezeHeaderRef={ref} metrics={bountyMetrics} />
          {loading ? (
            <LoaderContainer>
              <EuiLoadingSpinner size="l" />
            </LoaderContainer>
          ) : (
            <MyTable
              bounties={bounties}
              startDate={startDate}
              endDate={endDate}
              headerIsFrozen={inView}
              bountyStatus={bountyStatus}
              setBountyStatus={setBountyStatus}
              dropdownValue={dropdownValue}
              setDropdownValue={setDropdownValue}
            />
          )}
        </Container>
      )}
    </>
  );
};
