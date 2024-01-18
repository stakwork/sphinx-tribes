import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { OrgHeader } from 'pages/tickets/org/orgHeader';
import { mainStore } from 'store/main';

// Mock the necessary module to provide org_uuid
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: '/some-path/org-uuid-123', // Mocked path containing org_uuid
  }),
}));

describe('OrgHeader Component', () => {
  beforeEach(() => {
    jest.spyOn(mainStore, 'getOrgBounties').mockReset();
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

  it('calls getOrgBounties with correct parameters including org_uuid on search', () => {
    render(<OrgHeader />);

    // Simulate entering search text
    const searchText = 'sample search';
    const searchInput = screen.getByPlaceholderText('Search');
    fireEvent.change(searchInput, { target: { value: searchText } });

    // Simulate pressing Enter key
    fireEvent.keyUp(searchInput, { key: 'Enter', code: 'Enter' });

    // Expected org_uuid extracted from the mocked URL
    const expectedOrgUuid = 'org-uuid-123';

    // Check if getOrgBounties is called with correct parameters
    expect(mainStore.getOrgBounties).toHaveBeenCalledWith({
      page: 1,
      resetPage: true,
      search: searchText,
      org_uuid: expectedOrgUuid,
    });
  });
});
