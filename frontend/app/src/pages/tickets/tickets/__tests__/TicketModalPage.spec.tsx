import React from 'react';
import { render, act } from '@testing-library/react';
import '@testing-library/jest-dom';
import { Router, Route } from 'react-router-dom';
import { createMemoryHistory } from 'history';
import { TicketModalPage } from '../../TicketModalPage';

jest.mock('react-router-dom', () => {
  const originalModule = jest.requireActual('react-router-dom');
  return {
    ...originalModule,
    useHistory: () => ({
      push: jest.fn(),
      goBack: jest.fn()
    })
  };
});

describe('<TicketModalPage />', () => {
  it('should navigate to the correct URL when accessed directly and goBack is called', async () => {
    const history = createMemoryHistory({ initialEntries: ['/bounty/1181'] });
    const pushSpy = jest.spyOn(history, 'push');
    const goBackSpy = jest.spyOn(history, 'goBack');

    // Set up an empty referrer to simulate direct access
    Object.defineProperty(document, 'referrer', {
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
      expect(pushSpy).toHaveBeenCalledWith('/bounties');
      expect(goBackSpy).toHaveBeenCalled();
    });
  });
});
