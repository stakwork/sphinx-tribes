import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { TicketModalPage } from '../TicketModalPage';

// Mock required dependencies and store functions
jest.mock('hooks', () => ({
  useIsMobile: jest.fn(() => false), // Mocking useIsMobile as false
}));

jest.mock('store', () => ({
  useStores: jest.fn(() => ({
    main: {
      getBountyById: jest.fn(() => [{ person: { name: 'John' } }]), // Mocking main.getBountyById
      getBountyIndexById: jest.fn(() => 0), // Mocking main.getBountyIndexById
      peopleBounties: [], // Mocking main.peopleBounties
    },
    modals: {
      setStartupModal: jest.fn(), // Mocking modals.setStartupModal
    },
    ui: {
      meInfo: true, // Mocking ui.meInfo
    },
  })),
}));

describe('TicketModalPage', () => {
  it('renders modal content when visible', async () => {
    render(
      <BrowserRouter>
        <TicketModalPage setConnectPerson={() => {}} />
      </BrowserRouter>
    );

    expect(await screen.findByText('John')).toBeInTheDocument();
  });

  it('handles goBack function', () => {
    render(
      <BrowserRouter>
        <TicketModalPage setConnectPerson={() => {}} />
      </BrowserRouter>
    );

    fireEvent.click(screen.getByText('Go Back'));

    expect(screen.queryByText('John')).toBeNull();
  });

});
