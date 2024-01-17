import React from 'react';
import { render, act } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: jest.fn(),
    goBack: jest.fn()
  })
}));

describe('<TicketModalPage />', () => {
  it('should navigate to the correct URL when accessed directly and goBack is called', async () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });
    Object.defineProperty(window.document, 'referrer', {
      value: '',
      configurable: true
    });

    render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    await act(async () => {
      expect(history.push).toHaveBeenCalledWith('/bounties');
      // or expect(history.goBack).toHaveBeenCalled();
    });
  });
});
