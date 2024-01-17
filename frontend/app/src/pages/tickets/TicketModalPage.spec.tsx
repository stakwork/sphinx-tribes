import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { TicketModalPage } from './TicketModalPage.tsx'; // Adjust the import path as needed
import { BrowserRouter } from 'react-router-dom';

// Mocks
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: jest.fn(),
  useLocation: jest.fn(),
  useParams: jest.fn()
}));

const mockedHistoryPush = jest.fn();
const mockedUseHistory = require('react-router-dom').useHistory;
const mockedUseLocation = require('react-router-dom').useLocation;
const mockedUseParams = require('react-router-dom').useParams;

describe('TicketModalPage', () => {
  beforeEach(() => {
    mockedUseHistory.mockReturnValue({ push: mockedHistoryPush });
    mockedUseLocation.mockReturnValue({
      pathname: '/bounty/1203',
      search: ''
    });
    mockedUseParams.mockReturnValue({ uuid: 'ck1p7l6a5fdlqdgmmnpg', bountyId: '1203' });
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('redirects to home page on direct access to bounty page', async () => {
    jest.setTimeout(10000);
    // Arrange: Set up conditions for a direct access
    Object.defineProperty(window, 'referrer', {
      configurable: true,
      value: ''
    });

    // Act: Render the component and simulate closing the modal
    render(
      <BrowserRouter>
        <TicketModalPage setConnectPerson={() => {}} />
      </BrowserRouter>
    );
    const goBackButton = await waitFor(() => screen.findByTestId('testid-modal'), {
      timeout: 10000
    });

    fireEvent.click(goBackButton);

    // Assert: Expect to be redirected to the home page
    expect(mockedHistoryPush).toHaveBeenCalledWith('/bounties');
  });

  // Other tests for different scenarios...
});
