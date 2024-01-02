import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
import UserTickets from '../UserTicketsView.tsx';

const mockTickets = [
  {
    body: {
      id: 1,
      title: 'Mock Bounty Title 1',
      description: 'Mock Bounty Description 1',
      price: 1000,
      estimatedTime: '2 hours',
      owner_id: 'ownerId',
      created: new Date().toISOString()
    }
  },
  {
    body: {
      id: 2,
      title: 'Mock Bounty Title 2',
      description: 'Mock Bounty Description 2',
      price: 2000,
      estimatedTime: '3 hours',
      owner_id: 'ownerId',
      created: new Date().toISOString()
    }
  }
];

jest.mock('store', () => ({
  useStores: jest.fn(() => ({
    main: {
      getPersonAssignedBounties: jest.fn(() => Promise.resolve(mockTickets)),
      people: [{ id: 1, owner_pubkey: 'ownerId' }]
    },
    ui: {
      meInfo: {
        url: 'https://example.com',
        jwt: 'mock-jwt-token'
      },
      setBountyPerson: jest.fn()
    }
  }))
}));

describe('UserTickets Component', () => {
  test('renders UserTickets component and displays bounties', async () => {
    render(
      <Router>
        <UserTickets />
      </Router>
    );

    await screen.findByTestId('test');

    expect(screen.queryByText('No results found')).toBeNull();
    expect(screen.getAllByTestId('test').length).toBeGreaterThan(0);

    expect(screen.getByText(mockTickets[0].body.title)).toBeInTheDocument();

    expect(screen.getByText(mockTickets[0].body.description)).toBeInTheDocument();

    expect(screen.getByText(`${mockTickets[0].body.price} Sats`)).toBeInTheDocument();

    expect(screen.getByText(mockTickets[0].body.estimatedTime)).toBeInTheDocument();
  });
});
