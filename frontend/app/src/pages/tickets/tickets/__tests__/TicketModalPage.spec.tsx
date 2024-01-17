import React from 'react';
import { render } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage'; // Replace with your actual import

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
  it('should redirect to the home page on direct access', () => {
    const history = createMemoryHistory();
    jest.spyOn(document, 'referrer', 'get').mockReturnValue('');

    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    expect(history.location.pathname).toEqual('/org/bounties/ck9drb84nncjnaefo090');
  });

  it('should go back on non-direct access', () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });
    jest
      .spyOn(document, 'referrer', 'get')
      .mockReturnValue('https://community.sphinx.chat/bounties');

    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    expect(history.location.pathname).not.toEqual('/org/bounties/ck9drb84nncjnaefo090');
    expect(history.location.pathname).toEqual('https://community.sphinx.chat/bounties');
  });
});
