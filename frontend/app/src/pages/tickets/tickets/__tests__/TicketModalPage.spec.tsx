import React from 'react';
import { render, act } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage'; // Adjust import as necessary

// Mock the necessary hooks and modules
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: jest.fn(),
    goBack: jest.fn()
  }),
  useParams: () => ({
    uuid: 'ck9drb84nncjnaefo090',
    bountyId: '1186'
  }),
  useLocation: () => ({
    pathname: '/bounty/1186',
    search: '',
    state: null
  })
}));

describe('TicketModalPage Navigation', () => {
  it('should redirect to the home page on direct access', async () => {
    const history = createMemoryHistory();
    jest.spyOn(document, 'referrer', 'get').mockReturnValue('');

    await act(async () => {
      render(
        <Router history={history}>
          <Route path="/bounty/:bountyId">
            <TicketModalPage setConnectPerson={jest.fn()} />
          </Route>
        </Router>
      );
    });

    // Update your expectation as per the actual redirection logic
    expect(history.location.pathname).toEqual('/org/bounties/ck9drb84nncjnaefo090');
  });

  it('should go back on non-direct access', async () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });
    jest
      .spyOn(document, 'referrer', 'get')
      .mockReturnValue('https://community.sphinx.chat/bounties');

    await act(async () => {
      render(
        <Router history={history}>
          <Route path="/bounty/:bountyId">
            <TicketModalPage setConnectPerson={jest.fn()} />
          </Route>
        </Router>
      );
    });

    // Update your expectation as per the actual redirection logic
    expect(history.location.pathname).not.toEqual('/org/bounties/ck9drb84nncjnaefo090');
    expect(history.location.pathname).toEqual('/previous-path');
  });
});
