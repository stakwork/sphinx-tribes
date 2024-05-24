
export const HostName = 'localhost:5005';


export var User: Cypress.Person = {
    owner_pubkey: "test_pubkey",
    owner_alias: "Alice",
    unique_name: "Alice",
    description: "this is a test",
    tags: [],
    img: "",
    unlisted: false,
    deleted: false,
    owner_route_hint: "",
    owner_contact_key: "",
    price_to_meet: 0,
    twitter_confirmed: false,
    referred_by: 0,
    extras: {},
    github_issues: {}
}

export const Workspaces = [
    {
        id: 0,
        uuid: 'cohob00n1e4808utqel0',
        name: ' Sample Workspace 1 ',
        owner_pubkey: User.owner_pubkey,
        img: 'IqQnBnAdrteW_QCeq_3Ss1_78_yBAz_rckG5F3NE9ms=',
        mission: 'Access the largest pool of human cognition',
        tactics: 'Create a marketplace for providers and bounty hunters',
        schematic_url: 'https://miro.com/app/board/uXjVNQOK7Zc=',
        schematic_img: '1b281867-2c0e-481e-b508-1aab0e33ab50.jpg'
    },
    {
        id: 0,
        uuid: 'cohob3on1e4808utqelg',
        name: 'Sample Workspace 2',
        owner_pubkey: User.owner_pubkey,
        img: 'IqQnBnAdrteW_QCeq_3Ss1_78_yBAz_rckG5F3NE9ms=',
        mission: 'Sample mission for worspace 2',
        tactics: 'Sample tactics for workspace 2',
        schematic_url: 'https://miro.com/app/board/sampleworkspace2',
        schematic_img: '1b281867-2c0e-481e-b508-1aab0e33ab50.jpg'
    },
    {
        id: 0,
        uuid: 'cohob80n1e4808utqem0',
        name: 'Sample Workspace 3',
        owner_pubkey: User.owner_pubkey,
        img: 'IqQnBnAdrteW_QCeq_3Ss1_78_yBAz_rckG5F3NE9ms=',
        mission: 'Sample mission for workspaces 3',
        tactics: 'Sample tactics for workspace 3',
        schematic_url: 'https://miro.com/app/board/sampleworkspace3',
        schematic_img: '1b281867-2c0e-481e-b508-1aab0e33ab50.jpg'
    }];

export const Repositories = [
    {
        uuid: 'com1t3gn1e4a4qu3tnlg',
        workspace_uuid: 'cohob00n1e4808utqel0',
        name: ' frontend ',
        url: ' https://github.com/stakwork/sphinx-tribes-frontend '
    },
    {
        uuid: 'com1t3gn1e4a4qu3thss',
        workspace_uuid: 'cohob00n1e4808utqel0',
        name: ' backend ',
        url: ' https://github.com/stakwork/sphinx-tribes '
    }
];

export const Features = [
    {
        uuid: 'com1kson1e49th88dbg0',
        workspace_uuid: 'cohob00n1e4808utqel0',
        name: ' Hive Process ',
        priority: 1,
        brief: ' To follow a set of best practices in product development.</br>' +
            'Dividing complex features into small<br>steps makes it easier to ' +
            'track and the timing more certain.<br/>A guided process would help ' +
            'a PM new to the hive process get the best results with the least mental ' +
            'load.<br/>This feature is for a not se technical Product Manager.<br/>' +
            'The hive process lets you get features out to production faster and with less risk. ',
        requirements: ' Modify workspaces endpoint to accomodate new fields.<br/>' +
            'Create end points for features, user stories and phases ',
        architecture: ' Describe the architecture of the feature with the following sections:' +
            '<br/><br/>Wireframes<br/><br/>Visual Schematics<br/><br/>Object Definition<br/><br/>' +
            'DB Schema Changes<br/><br/>UX<br/><br/>CI/CD<br/><br/>Changes<br/><br/>Endpoints<br/><br/>' +
            'Front<br/><br/> ',
    },
    {
        uuid: 'com1l5on1e49tucv350g',
        workspace_uuid: 'cohob00n1e4808utqel0',
        name: ' AI Assited text fields ',
        priority: 2,
        brief: 'An important struggle of a technical product manager is to find ' +
            'the right words to describe a business goal. The definition of ' +
            'things like \'product mission\' or \'tactics and objectives\' is ' +
            ' the base from which every technical decition relays on.<br/>' +
            'We are going to leverage AI to help the PM write better definitions.<br/>' +
            'The fields that would benefit form AI assistance are: mission, tactics, ' +
            'feature brief and feature user stories ',
        requirements: ' Create a new page for a conversation format between the PM and the LLM<br/>' +
            'Rely as much as possible on stakwork workflows<br/>' +
            'Have history of previous definitions ',
        architecture: ' Describe the architecture of the feature with the following sections:' +
            '<br/><br/>Wireframes<br/><br/>Visual Schematics<br/><br/>Object Definition<br/><br/>' +
            'DB Schema Changes<br/><br/>UX<br/><br/>CI/CD<br/><br/>Changes<br/><br/>Endpoints<br/><br/>' +
            'Front<br/><br/> ',
    },
    {
        uuid: 'com1l5on1e49tucv350h',
        workspace_uuid: 'cohob00n1e4808utqel0',
        name: ' AI Assited relation between text fields ',
        priority: 2,
        brief: ' A product and feature\'s various definition fields: mission, tactics, ' +
            'feature brief, user stories, requirements and architecture should have some ' +
            'relation between each other.<br/>' + 'One way to do that is to leverage an LLM ' +
            'to discern the parts of the defintion that have a connection to other definitions.<br/>' +
            'The UI will need to show the user how each definition is related to other defintions. ',
        requirements: 'Create a new process after a Feature text has changed. It should use the LLM to ' +
            'determine de relationship between parts of the text. ',
        architecture: 'Describe the architecture of the feature with the following sections:' +
            '<br/><br/>Wireframes<br/><br/>Visual Schematics<br/><br/>Object Definition<br/><br/>' +
            'DB Schema Changes<br/><br/>UX<br/><br/>CI/CD<br/><br/>Changes<br/><br/>Endpoints<br/><br/>' +
            'Front<br/><br/> ',
    },
];

export const UserStories = [
    { uuid: 'com1lh0n1e49ug76noig', feature_uuid: 'com1kson1e49th88dbg0', description: ' As a {PM} I want to {make providers \"hive ready\"}, so I can {leverage the hive process ', priority: 0 },
    { uuid: 'com1lk8n1e49uqfe3l40', feature_uuid: 'com1kson1e49th88dbg0', description: ' As a {PM} I want to {CRUD Features}, so I can {use the system to manage my features} ', priority: 1 },
    { uuid: 'com1ln8n1e49v4159gug', feature_uuid: 'com1kson1e49th88dbg0', description: ' As a {PM} I want to {follow best practices}, so I can {make more valuable features} ', priority: 2 },
    { uuid: 'com1lqgn1e49vevhs9k0', feature_uuid: 'com1kson1e49th88dbg0', description: ' As a {PM} I want to {save the architecture of the feature}, so I can {share it with people} ', priority: 3 },
    { uuid: 'com1lt8n1e49voquoq90', feature_uuid: 'com1kson1e49th88dbg0', description: ' As a {PM} I want to {create phases}, so I can {divide the work in several deliverable stages} ', priority: 4 },
    { uuid: 'com1m08n1e4a02r6j0pg', feature_uuid: 'com1kson1e49th88dbg0', description: ' As a {PM} I want to {assign bounties to features}, so I can {group bounties together} ', priority: 5 },
];

export const Phases = [
    { uuid: 'com1msgn1e4a0ts5kls0', feature_uuid: 'com1l5on1e49tucv350g', name: ' MVP ', priority: 0 },
    { uuid: 'com1mvgn1e4a1879uiv0', feature_uuid: 'com1l5on1e49tucv350g', name: ' Phase 2 ', priority: 1 },
    { uuid: 'com1n2gn1e4a1i8p60p0', feature_uuid: 'com1l5on1e49tucv350g', name: ' Phase 3 ', priority: 2 },
];

export const Bounties = [
    {
        owner_id: 'alice',
        paid: false,
        show: true,
        completed: false,
        type: 'bug',
        award: '1000 USD',
        assigned_hours: 5,
        bounty_expires: '2024-06-01T00:00:00Z',
        commitment_fee: 1000,
        price: 21,
        title: 'Phase 1 Bounty 1',
        ticket_url: 'http://ticket.url',
        workspace_uuid: 'workspace-uuid-123',
        phase_uuid: Phases[1].uuid,
        phase_priority: 0,
        description: 'Detailed bug description',
        wanted_type: 'Bugfix',
        deliverables: 'Bug should be fixed',
        github_description: true,
        one_sentence_summary: 'Fix bug in production',
        estimated_session_length: '2 hours',
        estimated_completion_date: '2024-05-25T00:00:00Z',
        created: Date.now(),
        coding_languages: ['JavaScript', 'Python']
    },
    {
        id: 0,
        owner_id: 'alice',
        paid: false,
        show: true,
        completed: false,
        type: 'freelance_job_request',
        award: '', 
        assigned_hours: 2,
        bounty_expires: '2024-06-01T00:00:00Z',
        commitment_fee: 0,
        price: 21,
        title: 'Phase 1 Bounty 2',
        tribe: '',
        assignee: '',
        ticket_url: '',
        workspace_uuid: Workspaces[0].uuid,
        phase_uuid: Phases[2].uuid,
        phase_priority: 1,
        description: 'detailed bounty description',
        wanted_type: 'Web development',
        deliverables: '',
        github_description: true,
        one_sentence_summary: '',
        estimated_session_length: 'Less than 3 hours',
        estimated_completion_date: '2024-05-25T00:00:00Z',
        coding_languages: ['Golang'],
    },
    {
        id: 0,
        owner_id: 'alice',
        paid: false,
        show: true,
        completed: false,
        type: 'freelance_job_request',
        award: '', 
        assigned_hours: 2,
        bounty_expires: '2024-06-01T00:00:00Z',
        commitment_fee: 0,
        price: 21,
        title: 'Phase 1 Bounty 3',
        tribe: '',
        assignee: '',
        ticket_url: '',
        workspace_uuid: Workspaces[0].uuid,
        phase_uuid: Phases[1].uuid,
        phase_priority: 2,
        description: 'detailed bounty description',
        wanted_type: 'Web development',
        deliverables: '',
        github_description: true,
        one_sentence_summary: '',
        estimated_session_length: 'Less than 3 hours',
        estimated_completion_date: '2024-05-25T00:00:00Z',
        coding_languages: ['Golang'],
    },
    {
        id: 0,
        owner_id: 'alice',
        paid: false,
        show: true,
        completed: false,
        type: 'freelance_job_request',
        award: '', 
        assigned_hours: 2,
        bounty_expires: '2024-06-01T00:00:00Z',
        commitment_fee: 0,
        price: 21,
        title: 'Phase 1 Bounty 4',
        tribe: '',
        assignee: '',
        ticket_url: '',
        workspace_uuid: Workspaces[0].uuid,
        phase_uuid: Phases[1].uuid,
        phase_priority: 3,
        description: 'detailed bounty description',
        wanted_type: 'Web development',
        deliverables: '',
        github_description: true,
        one_sentence_summary: '',
        estimated_session_length: 'Less than 3 hours',
        estimated_completion_date: '2024-05-25T00:00:00Z',
        coding_languages: ['Golang'],
    },
    {
        id: 0,
        owner_id: 'alice',
        paid: false,
        show: true,
        completed: false,
        type: 'freelance_job_request',
        award: '', 
        assigned_hours: 2,
        bounty_expires: '2024-06-01T00:00:00Z',
        commitment_fee: 0,
        price: 21,
        title: 'Phase 1 Bounty 5',
        tribe: '',
        assignee: '',
        ticket_url: '',
        workspace_uuid: Workspaces[0].uuid,
        phase_uuid: Phases[2].uuid,
        phase_priority: 4,
        description: 'detailed bounty description',
        wanted_type: 'Web development',
        deliverables: '',
        github_description: true,
        one_sentence_summary: '',
        estimated_session_length: 'Less than 3 hours',
        estimated_completion_date: '2024-05-25T00:00:00Z',
        coding_languages: ['Golang'],
    },
    {
        id: 0,
        owner_id: 'alice',
        paid: false,
        show: true,
        completed: false,
        type: 'freelance_job_request',
        award: '', 
        assigned_hours: 2,
        bounty_expires: '2024-06-01T00:00:00Z',
        commitment_fee: 0,
        price: 21,
        title: 'Phase 2 Bounty 1',
        tribe: '',
        assignee: '',
        ticket_url: '',
        workspace_uuid: Workspaces[0].uuid,
        phase_uuid: Phases[1].uuid,
        phase_priority: 0,
        description: 'detailed bounty description',
        wanted_type: 'Web development',
        deliverables: '',
        github_description: true,
        one_sentence_summary: '',
        estimated_session_length: 'Less than 3 hours',
        estimated_completion_date: '2024-05-25T00:00:00Z',
        coding_languages: ['Golang'],
    },
    {
        id: 0,
        owner_id: 'alice',
        paid: false,
        show: true,
        completed: false,
        type: 'freelance_job_request',
        award: '', 
        assigned_hours: 2,
        bounty_expires: '2024-06-01T00:00:00Z',
        commitment_fee: 0,
        price: 21,
        title: 'Phase 2 Bounty 2',
        tribe: '',
        assignee: '',
        ticket_url: '',
        workspace_uuid: Workspaces[0].uuid,
        phase_uuid: Phases[1].uuid,
        phase_priority: 1,
        description: 'detailed bounty description',
        wanted_type: 'Web development',
        deliverables: '',
        github_description: true,
        one_sentence_summary: '',
        estimated_session_length: 'Less than 3 hours',
        estimated_completion_date: '2024-05-25T00:00:00Z',
        coding_languages: ['Golang'],
    },
]