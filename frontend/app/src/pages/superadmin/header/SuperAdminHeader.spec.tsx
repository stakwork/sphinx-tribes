import '@testing-library/jest-dom';
import { render, screen, fireEvent } from '@testing-library/react';
import moment from 'moment';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../__test__/__mockData__/setupStore';
import { user } from '../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../__test__/__mockFn__/useHistory';
import { Header } from './';

beforeAll(() => {
  nock.disableNetConnect();
  setupStore();
  mockUsehistory();
});

/**
 * @jest-environment jsdom
 */
describe('Header Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});

  test('display header with extras', () => {
    const setStartDateMock = jest.fn();
    const setEndDateMock = jest.fn();
    const hardCodedDateRange = '01 Oct - 31 Dec 2023';
    const exportCSVText = 'Export CSV';
    const initDateRange = '7 Days';

    render(
      <Header
        startDate={1703969172}
        endDate={1700073000}
        setStartDate={setStartDateMock}
        setEndDate={setEndDateMock}
      />
    );

    expect(screen.queryByText(hardCodedDateRange)).toBeInTheDocument();
    expect(screen.queryByText(exportCSVText)).toBeInTheDocument();
    expect(screen.queryByText(initDateRange)).toBeInTheDocument();

    fireEvent.click(screen.getByText('Last 7 Days'));

    expect(screen.getByText('7 Days')).toBeInTheDocument();
    expect(screen.getByText('30 Days')).toBeInTheDocument();
    expect(screen.getByText('90 Days')).toBeInTheDocument();
    expect(screen.getByText('Custom')).toBeInTheDocument();

    fireEvent.click(screen.getByText('30 Days'));

    const expectedStartDate = moment().subtract(30, 'days').startOf('day').unix();
    const expectedEndDate = moment().startOf('day').unix();

    expect(setStartDateMock).toHaveBeenCalledWith(expectedStartDate);
    expect(setEndDateMock).toHaveBeenCalledWith(expectedEndDate);
    expect(screen.queryByText(exportCSVText)).toBeInTheDocument();
    expect(screen.queryByText(initDateRange)).toBeInTheDocument();
  });
});
