import React from 'react';
import { render, fireEvent, waitFor, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { Router, Route } from 'react-router-dom';
import { TicketModalPage } from './TicketModalPage';
import { createMemoryHistory } from 'history';

describe('<TicketModalPage />', () => {
  it('should navigate to the correct URL when goBack is called via bigCloseImage', async () => {
    const history = createMemoryHistory();
    history.push('/bounty/1203');

    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    const closeButton = await waitFor(() => screen.getByTestId('close-btn'), { timeout: 5000 });

    fireEvent.click(closeButton);

    await waitFor(
      () => {
        expect(history.location.pathname).toBe('/bounties');
      },
      { timeout: 5000 }
    );
  });
});
