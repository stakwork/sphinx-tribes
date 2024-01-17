import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage';

describe('TicketModalPage', () => {
  it('redirects to the bounty home page on direct access and modal close', () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });
    const { getByRole } = render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    // Assuming there's a close button in your modal, you might need to adjust this selector
    const closeButton = getByRole('button', { name: /close/i });
    fireEvent.click(closeButton);

    // Check if the current URL is the bounty home page
    expect(history.location.pathname).toBe('/bounties');
  });
});
