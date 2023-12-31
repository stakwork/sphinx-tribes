import '@testing-library/jest-dom';
import { render, screen } from '@testing-library/react';
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
describe('AboutView Component', () => {
  nock(user.url).get('/person/id/1').reply(200, {});
  test('display about view with extras', () => {
    const setStartDateMock = jest.fn();
    const setEndDateMock = jest.fn();
    const hardCodedDateRange = '01 Oct - 31 Dec 2023';
    const exportCSVText = 'Export CSV';
    const initDateRange = '7 days';

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
  });
});
