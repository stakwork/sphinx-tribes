 const mockBounties = [
  {
    person: { name: "Satoshi Nakamoto", id: 1 },
    body: { content: "Detailed description of the bounty task" },
    org_uuid: "org123",
    title: "Fix Frontend Bug",
    description: "A bug in the front end needs fixing",
    owner_id: "owner001",
    created: 1670000000000,
    show: true,
    assignee: { name: "Jane Doe", id: 2 },
    wanted_type: "bugfix",
    type: "issue",
    price: "500",
    codingLanguage: "JavaScript",
    estimated_session_length: "2 hours",
    bounty_expires: "2024-01-01",
    commitment_fee: 50
  },
  {
    person: { name: "Alice Smith", id: 3 },
    body: { content: "Consectetur adipiscing elit..." },
    org_uuid: "org124",
    title: "Develop New Feature",
    description: "Development of a new feature for our app",
    owner_id: "owner002",
    created: 1670100000000,
    show: false,
    assignee: { name: "Bob Johnson", id: 4 },
    wanted_type: "feature",
    type: "enhancement",
    price: "1000",
    codingLanguage: "Python",
    estimated_session_length: "5 hours",
    bounty_expires: "2024-02-01",
    commitment_fee: 100
  },
  {
    person: { name: "Charlie Brown", id: 5 },
    body: { content: "Sed do eiusmod tempor incididunt..." },
    org_uuid: "org125",
    title: "Improve Database Performance",
    description: "Optimizing database for better performance",
    owner_id: "owner003",
    created: 1670200000000,
    show: true,
    assignee: { name: "Eve White", id: 6 },
    wanted_type: "optimization",
    type: "database",
    price: "750",
    codingLanguage: "SQL",
    estimated_session_length: "3 hours",
    bounty_expires: "2024-03-01",
    commitment_fee: 75
  },
  {
    person: { name: "Dave Green", id: 7},
    body: { content: "Ut enim ad minim veniam..." },
    org_uuid: "org126",
    title: "UI Redesign",
    description: "Complete overhaul of the user interface",
    owner_id: "owner004",
    created: 1670300000000,
    show: false,
    assignee: { name: "Frank Black", id: "a4" },
    wanted_type: "redesign",
    type: "ui",
    price: "1200",
    codingLanguage: "React",
    estimated_session_length: "6 hours",
    bounty_expires: "2024-04-01",
    commitment_fee: 120
  },
  {
    person: { name: "Ella Yellow", id: 8 },
    body: { content: "Quis nostrud exercitation ullamco laboris..." },
    org_uuid: "org127",
    title: "API Integration",
    description: "Integration of third-party APIs into our system",
    owner_id: "owner005",
    created: 1670400000000,
    show: true,
    assignee: { name: "Grace Violet", id: 9},
    wanted_type: "integration",
    type: "api",
    price: "900",
    codingLanguage: "Java",
    estimated_session_length: "4 hours",
    bounty_expires: "2024-05-01",
    commitment_fee: 90
  }
];

export default mockBounties;


export const newBounty = {
  person: { id: 10 , name: "Alice Johnson" },
  body: "Bounty details here",
    org_uuid: "org1",
    title: "Fix Backend Bug",
    description: "Fix a critical bug in the backend",
    owner_id: "owner123",
    created: 1610000000,
    show: true,
    assignee: { id: "dev123", name: "Jane Smith" },
    wanted_type: "bugfix",
    type: "backend",
    price: "1000",
    codingLanguage: "JavaScript",
    estimated_session_length: "3 hours",
    bounty_expires: "2023-12-31",
    commitment_fee: 50
};

export const updatedBounty = {
    person: { name: "Satoshi Nakamoto", id: 1 },
    body: { content: "Detailed description of the bounty task" },
    org_uuid: "org123",
    title: "Updated Title: Fix Backend bugs",
    description: "A bug in the front end needs fixing",
    owner_id: "owner001",
    created: 1670000000000,
    show: true,
    assignee: { name: "Jane Doe", id: "a1" },
    wanted_type: "bugfix",
    type: "issue",
    price: "500",
    codingLanguage: "JavaScript",
    estimated_session_length: "2 hours",
    bounty_expires: "2024-01-01",
    commitment_fee: 50
  }
