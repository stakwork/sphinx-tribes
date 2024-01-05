import React, { useState } from 'react';
import { render, fireEvent } from '@testing-library/react';
import { createMemoryHistory } from 'history';
import userEvent from '@testing-library/user-event';
import { Router } from 'react-router-dom';
import '@testing-library/jest-dom';
import { MyTable } from './index.tsx';
import { BountyStatus } from 'store/main.ts';

const mockBounties = [
  {
    id: 1,
    title: 'Bounty 1',
    date: '2023-01-01',
    dtgp: 100,
    assignee: 'Assignee 1',
    assigneeImage: 'assignee-image-1.jpg',
    provider: 'Provider 1',
    providerImage: 'provider-image-1.jpg',
    organization: 'Org 1',
    organizationImage: 'org-image-1.jpg',
    status: 'open'
  },
  {
    id: 2,
    title: 'Bounty 2',
    date: '2023-01-02',
    dtgp: 200,
    assignee: 'Assignee 2',
    assigneeImage: 'assignee-image-2.jpg',
    provider: 'Provider 2',
    providerImage: 'provider-image-2.jpg',
    organization: 'Org 2',
    organizationImage: 'org-image-2.jpg',
    status: 'assigned'
  },
  {
    id: 3,
    title: 'Bounty 3',
    date: '2023-01-03',
    dtgp: 300,
    assignee: 'Assignee 3',
    assigneeImage: 'assignee-image-3.jpg',
    provider: 'Provider 3',
    providerImage: 'provider-image-3.jpg',
    organization: 'Org 3',
    organizationImage: 'org-image-3.jpg',
    status: 'paid'
  }
];

type Bounty = {
  id: number;
  title: string;
  date: string;
  dtgp: number;
  assignee: string;
  assigneeImage: string;
  provider: string;
  providerImage: string;
  organization: string;
  organizationImage: string;
  status: string;
};

it('renders elements from TableProps in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText(mockBounties[0].title)).toBeInTheDocument();
});

it('renders "Sort By:" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('Sort By:')).toBeInTheDocument();
});

it('renders "Date" twice in the document', () => {
  const { getAllByText } = render(<MyTable bounties={mockBounties} />);
  expect(getAllByText('Date')).toHaveLength(2);
});

it('renders "Assignee" twice in the document', () => {
  const { getAllByText } = render(<MyTable bounties={mockBounties} />);
  expect(getAllByText('Assignee')).toHaveLength(2);
});

it('renders "Status" twice in the document', () => {
  const { getAllByText } = render(<MyTable bounties={mockBounties} />);
  expect(getAllByText('Status')).toHaveLength(2);
});

it('renders "Status:" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('Status:')).toBeInTheDocument();
});

it('renders "All" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('All')).toBeInTheDocument();
});

it('renders "Open" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('Open')).toBeInTheDocument();
});

it('renders "In Progress" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('In Progress')).toBeInTheDocument();
});

it('renders "Completed" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('Completed')).toBeInTheDocument();
});

it('renders "Bounty" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('Bounty')).toBeInTheDocument();
});

it('renders "#DTGP" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('#DTGP')).toBeInTheDocument();
});

it('renders "Provider" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('Provider')).toBeInTheDocument();
});

it('renders "Organization" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText('Organization')).toBeInTheDocument();
});

it('renders each element in the table in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  expect(getByText(mockBounties[0].title)).toBeInTheDocument();
});

it('renders each element in the table in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} />);
  mockBounties.forEach((bounty: Bounty) => {
    expect(getByText(bounty.title)).toBeInTheDocument();
    expect(getByText(bounty.date)).toBeInTheDocument();
    expect(getByText(String(bounty.dtgp))).toBeInTheDocument();
    expect(getByText(bounty.assignee)).toBeInTheDocument();
    expect(getByText(bounty.provider)).toBeInTheDocument();
    expect(getByText(bounty.status)).toBeInTheDocument();
    expect(getByText(bounty.organization)).toBeInTheDocument();
  });


  it('should navigate to the correct URL when a bounty is clicked', () => {
    const history = createMemoryHistory();
    const { getByText } = render(
      <Router history={history}>
        <MyTable bounties={mockBounties} />
      </Router>
    );
    const bountyTitle = getByText('Sample Bounty');
    fireEvent.click(bountyTitle);
    expect(history.location.pathname).toBe('/bounty/1');
  });

  it('renders correct color box for different bounty statuses', () => {
    const { getAllByTestId } = render(<MyTable bounties={mockBounties} />);
    const statusElements = getAllByTestId('bounty-status');
    expect(statusElements[0]).toHaveStyle('background-color: #618AFF');
    expect(statusElements[1]).toHaveStyle('background-color: #49C998');
    expect(statusElements[2]).toHaveStyle('background-color: #5F6368');
  });

it('it renders with filter status states', async () => {
  const [bountyStatus, setBountyStatus] = useState<BountyStatus>({
    Open: false,
    Assigned: false,
    Paid: false
  });
  const [dropdownValue, setDropdownValue] = useState('all');

  const { getByText, getByLabelText } = render(
    <MyTable
      bounties={mockBounties}
      dropdownValue={dropdownValue}
      setDropdownValue={setDropdownValue}
      bountyStatus={bountyStatus}
      setBountyStatus={setBountyStatus}
    />
  );

  const dropdown = getByLabelText('Status:');
  fireEvent.select(dropdown);
  await userEvent.click(getByText('Open'));
  expect(dropdownValue).toBe('open');
});
