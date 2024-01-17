import React from 'react';
import { render } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage'; // Adjust the import path as needed

describe('TicketModalPage', () => {
  it('redirects to the bounty home page on direct access and modal close', () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/123'] });

    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    history.goBack();

    expect(history.location.pathname).toBe('/bounties');
  });
});
