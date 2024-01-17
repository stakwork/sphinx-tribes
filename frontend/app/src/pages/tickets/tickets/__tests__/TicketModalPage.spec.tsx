import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage'; // Adjust the import path as needed

describe('TicketModalPage', () => {
  it('redirects to the bounty home page on direct access and modal close', () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });

    jest.mock('react-router-dom', () => ({
      ...jest.requireActual('react-router-dom'),
      useHistory: () => history
    }));

    const { getByTestId } = render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    // Simulate clicking the big close button in the modal
    const bigCloseButton = getByTestId('close-btn');
    fireEvent.click(bigCloseButton);

    // Check if the current URL is the bounty home page
    expect(history.location.pathname).toBe('/bounties');
  });
});
