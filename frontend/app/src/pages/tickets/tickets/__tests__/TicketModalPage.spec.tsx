import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { TicketModalPage } from '../../TicketModalPage';
import { BrowserRouter as Router } from 'react-router-dom';

describe('TicketModalPage Navigation Test', () => {
  it('should redirect to bounties home when accessed directly and closed', () => {
    // Mock `history.push` and `history.goBack`
    const mockPush = jest.fn();
    const mockGoBack = jest.fn();

    jest.mock('react-router-dom', () => ({
      ...jest.requireActual('react-router-dom'),
      useHistory: () => ({
        push: mockPush,
        goBack: mockGoBack
      })
    }));

    // Render the component
    const { getByTestId } = render(
      <Router>
        <TicketModalPage setConnectPerson={jest.fn()} />
      </Router>
    );

    // Simulate the user closing the modal
    fireEvent.click(getByTestId('close-btn'));

    // Check if the correct navigation function was called
    expect(mockPush).toHaveBeenCalledWith('/bounties');
    expect(mockGoBack).not.toHaveBeenCalled();
  });
});
