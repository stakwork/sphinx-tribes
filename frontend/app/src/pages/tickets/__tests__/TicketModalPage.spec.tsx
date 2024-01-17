import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../TicketModalPage.tsx';

jest.mock('hooks', () => ({
  useIsMobile: () => false
}));
jest.mock('store', () => ({}));

describe('<TicketModalPage />', () => {
  it('should navigate to the correct URL when accessed directly and goBack is called', () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });

    Object.defineProperty(document, 'referrer', {
      value: '',
      writable: true
    });

    const { getByTestId } = render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    const goBackButton = getByTestId('go-back-button');
    fireEvent.click(goBackButton);

    expect(history.location.pathname).toBe('/bounties');
  });
});
