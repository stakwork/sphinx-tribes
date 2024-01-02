import '@testing-library/jest-dom';
import { render, screen, fireEvent, within } from '@testing-library/react';
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
  beforeEach(() => {
    nock.cleanAll();
  });

  nock(user.url).get('/person/id/1').reply(200, {});

  test('displays header with extras', async () => {
    const setStartDateMock = jest.fn();
    const setEndDateMock = jest.fn();
    const exportCSVText = 'Export CSV';
    const initDateRange = 'Last 7 Days';

    render(
      <Header
        startDate={moment().subtract(7, 'days').startOf('day').unix()}
        endDate={moment().startOf('day').unix()}
        setStartDate={setStartDateMock}
        setEndDate={setEndDateMock}
      />
    );

    const today = moment().startOf('day');
    const expectedStartDate = today.clone().subtract(7, 'days');
    const expectedEndDate = today;

    const leftWrapperElement = screen.getByTestId('leftWrapper');
    const monthElement = within(leftWrapperElement).getByTestId('month');

    expect(monthElement).toBeInTheDocument();

    expect(monthElement).toHaveTextContent(
      `${expectedStartDate.format('DD-MMM')} - ${expectedEndDate.format('DD-MMM-YYYY')}`
    );

    expect(screen.getByText(exportCSVText)).toBeInTheDocument();
    expect(screen.getByText(initDateRange)).toBeInTheDocument();

    fireEvent.click(screen.getByText(initDateRange));

    await screen.findByText('30 Days');
    fireEvent.click(screen.getByText('Last 30 Days'));

    const expectedStartDate30DaysMode = today.clone().subtract(30, 'days');
    const expectedEndDate30DaysMode = today;

    expect(monthElement).toHaveTextContent(
      `${expectedStartDate30DaysMode.format('DD-MMM')} - ${expectedEndDate30DaysMode.format(
        'DD-MMM-YYYY'
      )}`
    );

    expect(screen.getByText('Last 30 Days')).toHaveClass('selected');
    fireEvent.click(screen.getByText('Last 30 Days'));

    await screen.findByText('90 Days');
    fireEvent.click(screen.getByText('Last 90 Days'));

    const expectedStartDate90DaysMode = today.clone().subtract(90, 'days');
    const expectedEndDate90DaysMode = today;

    expect(monthElement).toHaveTextContent(
      `${expectedStartDate90DaysMode.format('DD-MMM')} - ${expectedEndDate90DaysMode.format(
        'DD-MMM-YYYY'
      )}`
    );

    expect(screen.getByText('Last 90 Days')).toHaveClass('selected');
  });
});
