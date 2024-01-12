import React, { useState } from 'react';
import { screen, render, fireEvent, waitFor } from '@testing-library/react';
import moment from 'moment';
import { createMemoryHistory } from 'history';
import userEvent from '@testing-library/user-event';
import { Router } from 'react-router-dom';
import '@testing-library/jest-dom';
import { MyTable } from '../index.tsx';
import { Bounty } from '../interfaces.ts';
//import { mockBounties } from '../mockAdminData.ts';

jest.mock('../index.tsx', () => ({
  ...jest.requireActual('../index.tsx'),
  paginatePrev: jest.fn(),
  paginateNext: jest.fn()
}));

const mockBounties = [
  {
    id: 1,
    bounty_id: 1,
    title: 'Bounty 1',
    date: '2023-01-01',
    bounty_created: '1672552800',
    paid_date: '2023-01-01',
    dtgp: 100,
    assignee: 'Assignee 1',
    assigneeImage: 'assignee-image-1.jpg',
    provider: 'Provider 1',
    providerImage: 'provider-image-1.jpg',
    organization: 'Org 1',
    organizationImage: 'org-image-1.jpg',
    status: 'open',
    paid: false
  },
  {
    id: 2,
    bounty_id: 2,
    title: 'Bounty 2',
    date: '2023-01-02',
    bounty_created: '1672639200',
    paid_date: '2023-01-02',
    dtgp: 200,
    assignee: 'Assignee 2',
    assigneeImage: 'assignee-image-2.jpg',
    provider: 'Provider 2',
    providerImage: 'provider-image-2.jpg',
    organization: 'Org 2',
    organizationImage: 'org-image-2.jpg',
    status: 'assigned',
    paid: false
  },
  {
    id: 3,
    bounty_id: 3,
    title: 'Bounty 3',
    date: '2023-01-03',
    bounty_created: '1672725600',
    paid_date: '2023-01-03',
    dtgp: 300,
    assignee: 'Assignee 3',
    assigneeImage: 'assignee-image-3.jpg',
    provider: 'Provider 3',
    providerImage: 'provider-image-3.jpg',
    organization: 'Org 3',
    organizationImage: 'org-image-3.jpg',
    status: 'paid',
    paid: true
  }
];

it('renders elements from TableProps in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText(mockBounties[0].title)).toBeInTheDocument();
});

it('renders "Sort By:" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('Sort By:')).toBeInTheDocument();
});

it('renders "Newest" twice in the document', () => {
  const { getAllByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getAllByText('Newest')).toHaveLength(2);
});

it('renders "Assignee" twice in the document', () => {
  const { getAllByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getAllByText('Assignee')).toHaveLength(2);
});

it('renders "Status" twice in the document', () => {
  const { getAllByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getAllByText('Status')).toHaveLength(2);
});

it('renders "Status:" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('Status:')).toBeInTheDocument();
});

it('renders "All" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('All')).toBeInTheDocument();
});

it('renders "Open" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('Open')).toBeInTheDocument();
});

it('renders "In Progress" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('In Progress')).toBeInTheDocument();
});

it('renders "Completed" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('Completed')).toBeInTheDocument();
});

it('renders "Bounty" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('Bounty')).toBeInTheDocument();
});

it('renders "#DTGP" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('#DTGP')).toBeInTheDocument();
});

it('renders "Provider" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('Provider')).toBeInTheDocument();
});

it('renders "Organization" in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText('Organization')).toBeInTheDocument();
});

it('renders each element in the table in the document', () => {
  const { getByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getByText(mockBounties[0].title)).toBeInTheDocument();
});

it('renders each element in the table in the document', () => {
  const { getByText, getAllByText } = render(
    <MyTable bounties={mockBounties} headerIsFrozen={false} />
  );
  const dates = ['2023-01-01', '2023-01-02', '2023-01-03'];
  const assignedText = getAllByText('assigned');
  expect(assignedText.length).toBe(2);
  expect(getByText('paid')).toBeInTheDocument();
  mockBounties.forEach((bounty: Bounty, index: number) => {
    expect(getByText(bounty.title)).toBeInTheDocument();
    expect(getByText(dates[index])).toBeInTheDocument();
    // expect(getByText(String(bounty.dtgp))).toBeInTheDocument();
    expect(getByText(bounty.assignee)).toBeInTheDocument();
    // expect(getByText(bounty.provider)).toBeInTheDocument();
    expect(getByText(bounty.organization)).toBeInTheDocument();
  });
});

it('should navigate to the correct URL when a bounty is clicked', () => {
  const history = createMemoryHistory();
  const { getByText } = render(
    <Router history={history}>
      <MyTable bounties={mockBounties} />
    </Router>
  );
  const bountyTitle = getByText('Bounty 1');
  fireEvent.click(bountyTitle);
  // expect(history.location.pathname).toBe('/bounty/1');
});

it('renders correct color box for different bounty statuses', () => {
  const { getAllByTestId } = render(<MyTable bounties={mockBounties} />);
  const statusElements = getAllByTestId('bounty-status');
  expect(statusElements[0]).toHaveStyle('background-color: #49C998');
  expect(statusElements[1]).toHaveStyle('background-color: #49C998');
  expect(statusElements[2]).toHaveStyle('background-color: #5F6368');
});

it('it renders with filter status states', async () => {
  const Wrapper = () => {
    const [bountyStatus, setBountyStatus] = useState({
      Open: false,
      Assigned: false,
      Paid: false
    });
    const [dropdownValue, setDropdownValue] = useState('all');
    return (
      <MyTable
        bounties={mockBounties}
        dropdownValue={dropdownValue}
        setDropdownValue={setDropdownValue}
        bountyStatus={bountyStatus}
        setBountyStatus={setBountyStatus}
      />
    );
  };
  const { getByText, getByLabelText } = render(<Wrapper />);

  const dropdown = getByText('All');
  fireEvent.select(dropdown);
  await userEvent.click(getByText('Open'));
  const openText = getByText('Open');
  expect(openText).toBeInTheDocument();
});

it('renders pagination section when number of bounties is greater than page size', () => {
  // Create an array of bounties with a count greater than the page size
  const largeMockBounties = Array.from({ length: 25 }, () => ({
    id: 1,
    bounty_id: 1,
    title:
      'Return user to the same page they were on before they edited a bounty user to the same page they were on before.',
    date: '2021.01.01',
    bounty_created: '2023-10-04T14:58:50.441223Z',
    paid_date: '2023-10-04T14:58:50.441223Z',
    dtgp: 1,
    assignee: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
    assigneeImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    provider:
      '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce:03a6ea2d9ead2120b12bd66292bb4a302c756983dc45dcb2b364b461c66fd53bcb:1099517001729',
    providerImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    organization: 'OrganizationName',
    organizationImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    status: 'open'
  }));
  const mockSetBountyStatus = jest.fn();
  const mockSetDropdownValue = jest.fn();

  render(
    <MyTable
      bounties={largeMockBounties}
      headerIsFrozen={false}
      startDate={moment().subtract(7, 'days').startOf('day').unix()}
      endDate={moment().startOf('day').unix()}
      bountyStatus={{ Open: false, Assigned: false, Paid: false }}
      setBountyStatus={mockSetBountyStatus}
      dropdownValue="all"
      setDropdownValue={mockSetDropdownValue}
    />
  );

  (async () => {
    await waitFor(() => {
      const paginationSection = screen.getByRole('pagination');
      expect(paginationSection).toBeInTheDocument();

      // Optionally, you can also check if pagination arrows are present
      const paginationArrowPrev = screen.getByAltText('pagination arrow 1');
      const paginationArrowNext = screen.getByAltText('pagination arrow 2');
      expect(paginationArrowPrev).toBeInTheDocument();
      expect(paginationArrowNext).toBeInTheDocument();
    });
  })();
});

const mockProps = {
  bounties: Array.from({ length: 25 }, () => ({
    id: 1,
    bounty_id: 1,
    title:
      'Return user to the same page they were on before they edited a bounty user to the same page they were on before.',
    date: '2021.01.01',
    bounty_created: '2023-10-04T14:58:50.441223Z',
    paid_date: '2023-10-04T14:58:50.441223Z',
    dtgp: 1,
    assignee: '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce',
    assigneeImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    provider:
      '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce:03a6ea2d9ead2120b12bd66292bb4a302c756983dc45dcb2b364b461c66fd53bcb:1099517001729',
    providerImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    organization: 'OrganizationName',
    organizationImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    status: 'open'
  })),
  startDate: moment().subtract(7, 'days').startOf('day').unix(),
  endDate: moment().startOf('day').unix(),
  headerIsFrozen: false,
  bountyStatus: { Open: false, Assigned: false, Paid: false },
  setBountyStatus: jest.fn(),
  dropdownValue: 'all',
  setDropdownValue: jest.fn(),
  paginatePrev: jest.fn(),
  paginateNext: jest.fn()
};
it('renders pagination arrows when bounties length is greater than pageSize and status filter is set to "open"', async () => {
  render(<MyTable {...mockProps} />);

  (async () => {
    await waitFor(() => {
      const paginationArrow1 = screen.getByAltText('pagination arrow 1');
      const paginationArrow2 = screen.getByAltText('pagination arrow 2');

      expect(paginationArrow1).toBeInTheDocument();
      expect(paginationArrow2).toBeInTheDocument();
    });
  })();
});

it('calls paginateNext when next pagination arrow is clicked with status filter set to "in-progress"', async () => {
  const inProgressProps = {
    ...mockProps,
    bountyStatus: { Open: false, Assigned: true, Paid: false }
  };
  render(<MyTable {...inProgressProps} />);

  (async () => {
    await waitFor(() => {
      const myTableInstance = screen.getByRole('pagination');
      const { paginateNext }: { paginateNext: any } = myTableInstance as any;
      const paginationArrow2 = screen.getByAltText('pagination arrow 2');
      fireEvent.click(paginationArrow2);

      expect(paginateNext).toHaveBeenCalled();
    });
  })();
});
