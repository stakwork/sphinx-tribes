import { render, fireEvent } from '@testing-library/react';
import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { TicketModalPage } from './TicketModalPage';

describe('TicketModalPage', () => {
  test('calls prevArrHandler when the previous button is clicked', () => {
    const setConnectPerson = jest.fn();
    const history = { replace: jest.fn() };

    const { getByRole } = render(
      <Router>
        <TicketModalPage setConnectPerson={setConnectPerson} />
      </Router>
    );

    fireEvent.click(getByRole('button', { name: /previous/i }));

    expect(setConnectPerson).toHaveBeenCalled();
    expect(history.replace).toHaveBeenCalled();
  });
});

describe('TicketModalPage', () => {
  test('calls nextArrHandler when the next button is clicked', () => {
    const setConnectPerson = jest.fn();
    const history = { replace: jest.fn() };

    const { getByRole } = render(
      <Router>
        <TicketModalPage setConnectPerson={setConnectPerson} />
      </Router>
    );

    fireEvent.click(getByRole('button', { name: /next/i }));

    expect(setConnectPerson).toHaveBeenCalled();
    expect(history.replace).toHaveBeenCalled();
  });
});
