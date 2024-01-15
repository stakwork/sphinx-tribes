import '@testing-library/jest-dom';
import { render, screen, within, act } from '@testing-library/react';
import moment from 'moment';
import nock from 'nock';
import React from 'react';
import { setupStore } from '../../../../__test__/__mockData__/setupStore';
import { user } from '../../../../__test__/__mockData__/user';
import { mockUsehistory } from '../../../../__test__/__mockFn__/useHistory';
import { Header } from '..';

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

    const { rerender } = render(
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

    act(() => {
      rerender(
        <Header
          startDate={moment().subtract(30, 'days').startOf('day').unix()}
          endDate={moment().startOf('day').unix()}
          setStartDate={setStartDateMock}
          setEndDate={setEndDateMock}
        />
      );
    });

    const StartDate30 = today.clone().subtract(30, 'days');
    expect(monthElement).toHaveTextContent(
      `${StartDate30.format('DD-MMM')} - ${expectedEndDate.format('DD-MMM-YYYY')}`
    );

    act(() => {
      rerender(
        <Header
          startDate={moment().subtract(90, 'days').startOf('day').unix()}
          endDate={moment().startOf('day').unix()}
          setStartDate={setStartDateMock}
          setEndDate={setEndDateMock}
        />
      );
    });

    const StartDate90 = today.clone().subtract(90, 'days');
    expect(monthElement).toHaveTextContent(
      `${StartDate90.format('DD-MMM')} - ${expectedEndDate.format('DD-MMM-YYYY')}`
    );
  });
});
