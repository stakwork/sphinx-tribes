import React from 'react';
import { render, fireEvent, waitFor } from '@testing-library/react';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage.tsx';

describe('<TicketModalPage />', () => {
  it('should navigate to the correct URL when accessed directly and goBack is called', async () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });

    const { getByTestId } = render(
      <Router history={history}>
        <Route path="/bounty/:bountyId">
          <TicketModalPage setConnectPerson={jest.fn()} />
        </Route>
      </Router>
    );

    let bigCloseButton;
    await waitFor(() => {
      bigCloseButton = getByTestId('close-btn');
    });

    fireEvent.click(bigCloseButton);

    expect(history.location.pathname).toBe('/bounties');
  });
});
