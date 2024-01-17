import React from 'react';
import { render } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage';

describe('TicketModalPage', () => {
  it('redirects to the bounty home page on direct access', () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1186'] });
    jest.spyOn(history, 'push');

    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    expect(history.push).toHaveBeenCalledWith('/bounties');
  });
});
