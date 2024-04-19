
let hostname = 'localhost:5555';

let user = {
    pubkey: '',
    authToken: '',
    name: 'Alice'
};

let workspaces = [
    {
        id: '29a999fc-0424-46d2-ac04-ffc3e0635928',
        name: 'Sample Workspace 1',
        owner_pub_key: user.pubkey,
        img: 'IqQnBnAdrteW_QCeq_3Ss1_78_yBAz_rckG5F3NE9ms=',
        mission: 'Access the largest pool of human cognition',
        tactics:'Create a marketplace for providers and bounty hunters',
        schematic_url: 'https://miro.com/app/board/uXjVNQOK7Zc=',
        schematic_img: '1b281867-2c0e-481e-b508-1aab0e33ab50.jpg'
    },
    {
        id: 'e3e7396a-3308-48f4-acba-b3d3120b65f6',
        name: 'Sample Workspace 2',
        owner_pub_key: user.pubkey,
        img: 'IqQnBnAdrteW_QCeq_3Ss1_78_yBAz_rckG5F3NE9ms=',
        mission: 'Access the largest pool of human cognition',
        tactics:'Create a marketplace for providers and bounty hunters',
        schematic_url: 'https://miro.com/app/board/uXjVNQOK7Zc=',
        schematic_img: '1b281867-2c0e-481e-b508-1aab0e33ab50.jpg'
    },
    {
        id: '18616ed4-2fda-4670-9550-dfa9356a3beb',
        name: 'Sample Workspace 3',
        owner_pub_key: user.pubkey,
        img: 'IqQnBnAdrteW_QCeq_3Ss1_78_yBAz_rckG5F3NE9ms=',
        mission: 'Access the largest pool of human cognition',
        tactics:'Create a marketplace for providers and bounty hunters',
        schematic_url: 'https://miro.com/app/board/uXjVNQOK7Zc=',
        schematic_img: '1b281867-2c0e-481e-b508-1aab0e33ab50.jpg'
    }];

let repositories = [
    {
        name: 'frontend',
        url: 'https://github.com/stakwork/sphinx-tribes-frontend'
    },
    {
        name: 'backend',
        url: 'https://github.com/stakwork/sphinx-tribes'
    }
];

let features = [
    {
        name: 'Hive Process',
        priority: 1,
        brief: 'To follow a set of best practices in product development.</br>' +
               'Dividing complex features into small<br>steps makes it easier to ' +
               'track and the timing more certain.<br/>A guided process would help ' +
               'a PM new to the hive process get the best results with the least mental ' +
               'load.<br/>This feature is for a not se technical Product Manager.<br/>' +
               'The hive process lets you get features out to production faster and with less risk.',
        requirements: 'Modify workspaces endpoint to accomodate new fields.<br/>' +
                      'Create end points for features, user stories and phases',
        architecture: 'Describe the architecture of the feature with the following sections:' +
                      '<br/><br/>Wireframes<br/><br/>Visual Schematics<br/><br/>Object Definition<br/><br/>' +
                      'DB Schema Changes<br/><br/>UX<br/><br/>CI/CD<br/><br/>Changes<br/><br/>Endpoints<br/><br/>' +
                      'Front<br/><br/>',
    },
    {
        name: 'AI Assited text fields',
        priority: 2,
        brief: 'An important struggle of a technical product manager is to find ' +
               'the right words to describe a business goal. The definition of ' +
               'things like \'product mission\' or \'tactics and objectives\' is ' +
               ' the base from which every technical decition relays on.<br/>' +
               'We are going to leverage AI to help the PM write better definitions.<br/>' +
               'The fields that would benefit form AI assistance are: mission, tactics, ' +
               'feature brief and feature user stories',
        requirements: 'Create a new page for a conversation format between the PM and the LLM<br/>' +
                      'Rely as much as possible on stakwork workflows<br/>' +
                      'Have history of previous definitions',
        architecture: 'Describe the architecture of the feature with the following sections:' +
                      '<br/><br/>Wireframes<br/><br/>Visual Schematics<br/><br/>Object Definition<br/><br/>' +
                      'DB Schema Changes<br/><br/>UX<br/><br/>CI/CD<br/><br/>Changes<br/><br/>Endpoints<br/><br/>' +
                      'Front<br/><br/>',
    },
    {
        name: 'AI Assited relation between text fields',
        priority: 2,
        brief: 'A product and feature\'s various definition fields: mission, tactics, ' +
               'feature brief, user stories, requirements and architecture should have some ' +
               'relation between each other.<br/>' + 'One way to do that is to leverage an LLM ' +
               'to discern the parts of the defintion that have a connection to other definitions.<br/>' +
               'The UI will need to show the user how each definition is related to other defintions.',
        requirements: 'Create a new process after a Feature text has changed. It should use the LLM to ' +
                      'determine de relationship between parts of the text.',
        architecture: 'Describe the architecture of the feature with the following sections:' +
                      '<br/><br/>Wireframes<br/><br/>Visual Schematics<br/><br/>Object Definition<br/><br/>' +
                      'DB Schema Changes<br/><br/>UX<br/><br/>CI/CD<br/><br/>Changes<br/><br/>Endpoints<br/><br/>' +
                      'Front<br/><br/>',
    },
];

let user_stories = [
    { id: 'f4c4c4b4-7a90-4a3a-b3e2-151d0feca9bf', description: ' As a {PM} I want to {make providers \"hive ready\"}, so I can {leverage the hive process ' },
    { id: '78f4b326-1841-449b-809a-a0947622db3e', description: ' As a {PM} I want to {CRUD Features}, so I can {use the system to manage my features} ' }, 
    { id: '5d353d23-3d27-4aa8-a9f7-04dcd5f4843c', description: ' As a {PM} I want to {follow best practices}, so I can {make more valuable features} ' },
    { id: '1a4e00f4-0e58-4e08-a1df-b623bc10f08d', description: ' As a {PM} I want to {save the architecture of the feature}, so I can {share it with people} ' }, 
    { id: 'eb6e4138-37e5-465d-934e-18e335abaa47', description: ' As a {PM} I want to {create phases}, so I can {divide the work in several deliverable stages} ' },
    { id: '35a5d8dd-240d-4ff0-a699-aa2fa2cfa32c', description: ' As a {PM} I want to {assign bounties to features}, so I can {group bounties together} ' },
];  

let phases = [
    { id: 'a96e3bff-e5c8-429e-bd65-911d619761aa', name: ' MVP ' },
    { id: '6de147ab-695c-45b1-81e7-2d1a5ba482ab', name: ' MVP ' },
    { id: '28541c4a-41de-447e-86d8-293583d1abc2', name: ' MVP ' },
];