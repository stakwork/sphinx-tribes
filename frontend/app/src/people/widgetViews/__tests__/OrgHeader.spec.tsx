import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { OrgHeader } from 'pages/tickets/org/orgHeader';
import { mainStore } from 'store/main';

describe('OrgHeader Component', () => {
  beforeEach(() => {
    jest.spyOn(mainStore, 'getPeopleBounties').mockReset();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

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
  });

  it('displays the correct number of bounties', () => {
    render(<OrgHeader />);
    expect(screen.getByText('284')).toBeInTheDocument();
    expect(screen.getByText('Bounties')).toBeInTheDocument();
  });

  it('should calls getPeopleBounties with correct parameters', () => {
    const orgUuid = 'cmkln4tm098m49vhlt80';
    Object.defineProperty(window, 'location', {
      value: {
        pathname: `/org/bounties/${orgUuid}`
      },
      writable: true
    });

    render(<OrgHeader />);

    // Simulate entering search text
    const searchText = 'sample search';
    const searchInput = screen.getByPlaceholderText('Search') as HTMLInputElement;
    fireEvent.change(searchInput, { target: { value: searchText } });

    // Simulate pressing Enter key
    fireEvent.keyUp(searchInput, { key: 'Enter', code: 'Enter' });

    // Check if getPeopleBounties is called with correct parameters
    expect(mainStore.getPeopleBounties).toHaveBeenCalledWith({
      page: 1,
      resetPage: true,
      search: searchText,
      org_uuid: orgUuid
    });
  });
});
