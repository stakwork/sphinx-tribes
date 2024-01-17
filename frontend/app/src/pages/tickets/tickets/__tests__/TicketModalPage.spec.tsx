import React from 'react';
import { render, fireEvent, waitFor, screen } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage.tsx';

describe('<TicketModalPage />', () => {
  it('should navigate to the correct URL when accessed directly and goBack is called', async () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });

    // Mock referrer if necessary
    Object.defineProperty(document, 'referrer', {
      value: '',
      writable: true
    });

    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    // Wait for the modal to be in the document
    const goBackButton = await waitFor(() => screen.findByTestId('close-btn'));

    fireEvent.click(goBackButton);

    expect(history.location.pathname).toBe('/bounties');
  });
});
