import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';
import { MyTable } from './index.tsx';

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
  }
];

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
  mockBounties.forEach((bounty) => {
    expect(getByText(bounty.title)).toBeInTheDocument();
    expect(getByText(bounty.date)).toBeInTheDocument();
    expect(getByText(String(bounty.dtgp))).toBeInTheDocument();
    expect(getByText(bounty.assignee)).toBeInTheDocument();
    expect(getByText(bounty.provider)).toBeInTheDocument();
    expect(getByText(bounty.organization)).toBeInTheDocument();
  });
});
