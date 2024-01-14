import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { OrgHeader } from 'pages/tickets/org/orgHeader';

describe('OrgHeader Component', () => {
  it('renders the component correctly', () => {
    render(<OrgHeader />);
    expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
    expect(screen.getByLabelText('Status')).toBeInTheDocument();
    expect(screen.getByLabelText('Skill')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument();
    expect(screen.getByText(/Bounties/i)).toBeInTheDocument();
  });

  it('opens the PostModal on "Post a Bounty" button click', async () => {
    render(<OrgHeader />);
    fireEvent.click(screen.getByText('Post a Bounty'));
    const modalTitle = await screen.findByText('Choose Bounty type', {}, { timeout: 1000 });
    expect(modalTitle).toBeInTheDocument();
  });

  it('displays the correct number of bounties', () => {
    render(<OrgHeader />);
    expect(screen.getByText('284')).toBeInTheDocument();
    expect(screen.getByText('Bounties')).toBeInTheDocument();
  });
});
