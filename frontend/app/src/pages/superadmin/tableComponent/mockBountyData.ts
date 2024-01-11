interface Bounty {
  id: number;
  bounty_id: number;
  title: string;
  date: string;
  bounty_created: string;
  paid_date: string;
  dtgp: number;
  assignee: string;
  assigneeImage: string;
  provider: string;
  providerImage: string;
  organization: string;
  organizationImage: string;
  status: string;
}

export const bounties: Bounty[] = [
  {
    id: 1,
    bounty_id: 1,
    title:
      'Return user to the same page they were on before they edited a bounty user to the same page they were on before.',
    date: '2021.01.01',
    bounty_created: '2021.01.01',
    paid_date: '',
    dtgp: 1,
    assignee:
      '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce:03a6ea2d9ead2120b12bd66292bb4a302c756983dc45dcb2b364b461c66fd53bcb:1099517001729',
    assigneeImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    provider:
      '035f22835fbf55cf4e6823447c63df74012d1d587ed60ef7cbfa3e430278c44cce:03a6ea2d9ead2120b12bd66292bb4a302c756983dc45dcb2b364b461c66fd53bcb:1099517001729',
    providerImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    organization: 'OrganizationName',
    organizationImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    status: 'Open',
  },
  {
    id: 2,
    bounty_id: 2,
    title:
      'Create a new website for the company Doe Inc. in React. The website should be responsive and have a dark mode.',
    date: '2021.01.01',
    bounty_created: '2021.01.01',
    paid_date: '',
    dtgp: 1,
    assignee: 'John Doe',
    assigneeImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    provider: 'Jane Doe',
    providerImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    organization: 'Doe Inc.',
    organizationImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    status: 'Paid',
  },
  {
    id: 3,
    bounty_id: 3,
    title:
      'Create a new website for the company Doe Inc. in React. The website should be responsive and have a dark mode.',
    date: '2023.08.05',
    bounty_created: '2023.08.05',
    paid_date: '',
    dtgp: 1,
    assignee: 'John Doe',
    assigneeImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    provider: 'Jane Doe',
    providerImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    organization: 'Doe Inc.',
    organizationImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    status: 'Assigned',
  },
  {
    id: 4,
    bounty_id: 4,
    title:
      'Create a new website for the company Doe Inc. in React. The website should be responsive and have a dark mode.',
    date: '2023.08.05',
    bounty_created: '2023.08.05',
    paid_date: '',
    dtgp: 1,
    assignee: 'John Doe',
    assigneeImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    provider: 'Jane Doe',
    providerImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    organization: 'Doe Inc.',
    organizationImage:
      'https://avatars.githubusercontent.com/u/10001?s=460&u=8c61f1cda5e9e2c2d1d5b8d2a5a8a5b8d2a5a8a5&v=4',
    status: 'Completed',
  },
];
