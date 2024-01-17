import React from 'react';
import { render, waitFor } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage';

// Mock useHistory
const mockHistoryPush = jest.fn();
const mockHistoryGoBack = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockHistoryPush,
    goBack: mockHistoryGoBack
  })
}));

describe('<TicketModalPage />', () => {
  beforeEach(() => {
    // Clear all mock function calls before each test
    mockHistoryPush.mockClear();
    mockHistoryGoBack.mockClear();
  });

  it('should navigate to the correct URL when accessed directly and goBack is called', async () => {
    // Create a memory history for your test
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });

    // Render the component inside a Router context
    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    await waitFor(() => {
      expect(mockHistoryPush).toHaveBeenCalledWith('/bounties');
      // or expect(mockHistoryGoBack).toHaveBeenCalled();
    });
  });
});
