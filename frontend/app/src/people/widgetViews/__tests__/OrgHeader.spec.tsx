import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MemoryRouter } from 'react-router-dom';
import { OrgHeader } from 'pages/tickets/org/orgHeader';

describe('OrgHeader Component', () => {
  it('renders the component correctly', () => {
    render(
      <MemoryRouter>
        <OrgHeader />
      </MemoryRouter>
    );
    expect(screen.getByText('Post a Bounty')).toBeInTheDocument();
    expect(screen.getByLabelText('Status')).toBeInTheDocument();
    expect(screen.getByLabelText('Skill')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Search')).toBeInTheDocument();
    expect(screen.getByText(/Bounties/i)).toBeInTheDocument();
  });

  it('opens the PostModal on "Post a Bounty" button click', async () => {
    render(
      <MemoryRouter>
        <OrgHeader />
      </MemoryRouter>
    );
    fireEvent.click(screen.getByText('Post a Bounty'));
  });

  it('displays the correct number of bounties', () => {
    render(
      <MemoryRouter>
        <OrgHeader />
      </MemoryRouter>
    );
    expect(screen.getByText('284')).toBeInTheDocument();
    expect(screen.getByText('Bounties')).toBeInTheDocument();
  });
});
