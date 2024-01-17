import React from 'react';
import { render, screen } from '@testing-library/react';
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
    render(
      <Router>
        <TicketModalPage setConnectPerson={jest.fn()} />
      </Router>
    );

    screen.debug();
    expect(screen.getByTestId('close-btn')).toBeInTheDocument();

    // Check if the correct navigation function was called
    expect(mockPush).toHaveBeenCalledWith('/bounties');
    expect(mockGoBack).not.toHaveBeenCalled();
  });
});
