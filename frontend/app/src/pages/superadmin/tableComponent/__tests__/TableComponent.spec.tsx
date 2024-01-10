import React, { useState } from 'react';
import { screen, render, fireEvent } from '@testing-library/react';
import { createMemoryHistory } from 'history';
import userEvent from '@testing-library/user-event';
import { Router } from 'react-router-dom';
import '@testing-library/jest-dom';
import { MyTable } from '../index.tsx';
import { BountyStatus } from 'store/main.ts';
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

it('renders "Date" twice in the document', () => {
  const { getAllByText } = render(<MyTable bounties={mockBounties} headerIsFrozen={false} />);
  expect(getAllByText('Date')).toHaveLength(2);
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
describe('MyTable Pagination', () => {
  // const mockBounties: TableProps['bounties'] = [
  //   // Your mock data here
  // ];

  const mockbounties = new Array(30).fill({
    owner_id: '021ae436bcd40ca21396e59be8cdb5a707ceacdb35c1d2c5f23be7584cab29c40b',
    paid: false,
    show: true,
    type: 'freelance_job_request',
    award: '',
    assigned_hours: 0,
    bounty_expires: '',
    commitment_fee: 0,
    price: 500000,
    title: 'Auto-populate org when "create a bounty" modal is initiated from org page',
    tribe: 'None',
    assignee: '02fcb9c3fb7b754acbd4755c2ee3fbeb6525eecb7418b19da96c14616158eae77d',
    ticket_url: 'https://github.com/stakwork/sphinx-tribes/issues/1320',
    org_uuid: 'ck95pe04nncjnaefo08g',
    wanted_type: 'Web development',
    deliverables: '',
    github_description: true,
    one_sentence_summary: '',
    estimated_session_length: '',
    estimated_completion_date: '',
    assigned_date: '2024-01-10T05:33:39.631167Z',
    coding_languages: [],
    bounty_id: 1132,
    bounty_created: 1704844695,
    bounty_updated: '2024-01-10T05:33:39.631167Z',
    bounty_description:
      '### Context\r\nWhen someone creates a bounty from an org page, the org field should be auto-populated when they reach step 2 instead of it being manually set. Here is an example of an org page: https://community.sphinx.chat/org/bounties/ck95pe04nncjnaefo08g\r\n\r\n### Design\r\n<img width="925" alt="Screenshot 2024-01-09 at 3 55 12 PM" src="https://github.com/stakwork/sphinx-tribes/assets/32662508/3569d1b2-5b77-4821-bea0-c4222e84b5d8">\r\n\r\n### Acceptance Criteria\r\n- [ ] I\'ve tested on Chrome\r\n- [ ] I\'ve submitted a recording in my pr showing this work\r\n- [ ] I\'ve included appropriate tests that asserts the same org is populated as the org page I\'m on',
    uuid: 'cjrptro4nncma8v11kmg',
    owner_pubkey: '021ae436bcd40ca21396e59be8cdb5a707ceacdb35c1d2c5f23be7584cab29c40b',
    unique_name: 'ecurrencyhodler',
    tags: [],
    img: 'https://memes.sphinx.chat/public/SlXI397USVtZpdltwXaA0yb9IuEnF-UVY-3QRwpWq1w=',
    unlisted: false,
    deleted: false,
    last_login: 1704834954,
    price_to_meet: 0,
    new_ticket_time: 0,
    twitter_confirmed: false,
    extras: {
      alert: false,
      amboss: [{ label: '', value: '' }],
      coding_languages: [
        {
          background: 'rgba(184, 37, 95, 0.1)',
          border: '1px solid rgba(184, 37, 95, 0.1)',
          color: '#B8255F',
          label: 'Lightning',
          value: 'Lightning'
        }
      ],
      email: [{ label: '', value: '' }],
      facebook: [{ label: '', value: '' }],
      github: [{ label: 'ecurrencyhodler', value: 'ecurrencyhodler' }],
      lightning: [{ label: 'ecurrencyhodler@getalby.com', value: 'ecurrencyhodler@getalby.com' }],
      tribes: [],
      twitter: [{ label: '', value: '' }]
    },
    github_issues: {},
    assignee_alias: 'Vayras',
    assignee_id: 243,
    assignee_img: 'https://memes.sphinx.chat/public/VvAAScl6-U2v8bt8roYuHbBPy2h4pb_t-G5H-fSsP3g=',
    assignee_created: null,
    assignee_updated: null,
    assignee_description: '',
    assignee_route_hint:
      '03a6ea2d9ead2120b12bd66292bb4a302c756983dc45dcb2b364b461c66fd53bcb:1099519164417',
    bounty_owner_id: 180,
    owner_uuid: 'cjrptro4nncma8v11kmg',
    owner_key: '',
    owner_alias: '',
    owner_unique_name: 'ecurrencyhodler',
    owner_description: 'Bitcoin PM',
    owner_tags: null,
    owner_img: 'https://memes.sphinx.chat/public/SlXI397USVtZpdltwXaA0yb9IuEnF-UVY-3QRwpWq1w=',
    owner_created: null,
    owner_updated: null,
    owner_last_login: 0,
    owner_route_hint: '',
    owner_contact_key: '',
    owner_price_to_meet: 0,
    owner_twitter_confirmed: false,
    organization_name: 'Bounties Platform',
    organization_img:
      'https://memes.sphinx.chat/public/IqQnBnAdrteW_QCeq_3Ss1_78_yBAz_rckG5F3NE9ms=',
    organization_uuid: 'ck95pe04nncjnaefo08g'
  });

  it('renders pagination arrows when bounties length is greater than pageSize', () => {
    const Wrapper = () => {
      const [bountyStatus, setBountyStatus] = useState({
        Open: false,
        Assigned: false,
        Paid: false
      });
      const [dropdownValue, setDropdownValue] = useState('all');
      const mockProps = {
        bounties: mockbounties,
        startDate: 1234567890,
        endDate: 1234567890,
        headerIsFrozen: false,
        bountyStatus: bountyStatus,
        setBountyStatus: setBountyStatus,
        dropdownValue: dropdownValue,
        setDropdownValue: setDropdownValue,
        paginatePrev: jest.fn(),
        paginateNext: jest.fn()
      };
      return <MyTable {...mockProps} />;
    };

    render(<Wrapper />);

    const paginationArrow1 = screen.getByAltText('pagination arrow 1');
    const paginationArrow2 = screen.getByAltText('pagination arrow 2');

    expect(paginationArrow1).toBeInTheDocument();
    expect(paginationArrow2).toBeInTheDocument();
  });

  it('calls paginateNext when next pagination arrow is clicked', () => {
    const Wrapper = () => {
      const [bountyStatus, setBountyStatus] = useState({
        Open: false,
        Assigned: false,
        Paid: false
      });
      const [dropdownValue, setDropdownValue] = useState('all');
      const mockProps = {
        bounties: mockbounties,
        startDate: 1234567890,
        endDate: 1234567890,
        headerIsFrozen: false,
        bountyStatus: bountyStatus,
        setBountyStatus: setBountyStatus,
        dropdownValue: dropdownValue,
        setDropdownValue: setDropdownValue,
        paginatePrev: jest.fn(),
        paginateNext: jest.fn()
      };
      return <MyTable {...mockProps} />;
    };

    render(<Wrapper />);
    const myTableInstance = screen.getByRole('pagination'); // Assuming role attribute is set appropriately
    const { paginateNext }: { paginateNext: any } = myTableInstance as any;
    const paginationArrow2 = screen.getByAltText('pagination arrow 2');
    fireEvent.click(paginationArrow2);

    expect(paginateNext).toHaveBeenCalled();
  });

  it('calls paginatePrev when previous pagination arrow is clicked', () => {
    const Wrapper = () => {
      const [bountyStatus, setBountyStatus] = useState({
        Open: false,
        Assigned: false,
        Paid: false
      });
      const [dropdownValue, setDropdownValue] = useState('all');
      const mockProps = {
        bounties: mockbounties,
        startDate: 1234567890,
        endDate: 1234567890,
        headerIsFrozen: false,
        bountyStatus: bountyStatus,
        setBountyStatus: setBountyStatus,
        dropdownValue: dropdownValue,
        setDropdownValue: setDropdownValue,
        paginatePrev: jest.fn(),
        paginateNext: jest.fn()
      };
      return <MyTable {...mockProps} />;
    };
    render(<Wrapper />);

    const myTableInstance = screen.getByRole('pagination'); // Assuming role attribute is set appropriately
    const { paginatePrev }: { paginatePrev: any } = myTableInstance as any;
    const paginationArrow1 = screen.getByAltText('pagination arrow 1');
    fireEvent.click(paginationArrow1);

    expect(paginatePrev).toHaveBeenCalled();
  });
});
