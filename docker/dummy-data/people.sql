-- Create people records for bounty developers
INSERT INTO people (
    id, 
    owner_pub_key, 
    owner_alias, 
    unique_name,
    description,
    img
) VALUES
-- Developer 1: High performer with many bounties
(
    1, 
    'developer_1', 
    'Sarah',
    'sarahj',
    'Senior backend developer specializing in Go and distributed systems. Previously at Stripe and Cloudflare.',
    'https://randomuser.me/api/portraits/women/22.jpg'
),

-- Developer 2: Moderate performer with specialized skills
(
    2, 
    'developer_2', 
    'Chen',
    'mikechen',
    'Full-stack developer with focus on real-time systems and WebSockets. I love building reliable, scalable applications.',
    'https://randomuser.me/api/portraits/men/34.jpg'
),

-- Developer 3: Highest earner with fewer but expensive bounties
(
    3, 
    'developer_3', 
    'Elena',
    'elenarodriguez',
    'Blockchain expert specializing in Bitcoin and smart contracts. Contributor to multiple open-source crypto projects.',
    'https://randomuser.me/api/portraits/women/45.jpg'
),

-- Developer 4: New but promising
(
    4, 
    'developer_4', 
    'Alex Thompson',
    'alexthompson',
    'Frontend developer specializing in React and modern UI/UX. Passionate about creating beautiful, accessible interfaces.',
    'https://randomuser.me/api/portraits/men/53.jpg'
),

-- Developer 5: Specialist in docs and testing
(
    5, 
    'developer_5', 
    'Priya Patel',
    'priyapatel',
    'Technical writer and documentation specialist. Converting complex topics into clear, concise documentation.',
    'https://randomuser.me/api/portraits/women/67.jpg'
),

-- Developer 6: New addition to the team
(
    6, 
    'developer_6', 
    'James Wilson',
    'jameswilson',
    'Full-stack developer with expertise in frontend frameworks and admin dashboards. Former designer turned coder.',
    'https://randomuser.me/api/portraits/men/75.jpg'
),

-- Developer 7: DevOps specialist
(
    7, 
    'developer_7', 
    'David Kim',
    'davidkim',
    'DevOps engineer specializing in CI/CD pipelines and infrastructure automation. Keeping systems reliable and scalable.',
    'https://randomuser.me/api/portraits/men/82.jpg'
),

-- Developer 8: Backend systems architect
(
    8, 
    'developer_8', 
    'Olivia Martinez',
    'oliviamartinez',
    'Backend systems architect with focus on high-performance distributed systems. I build things that scale.',
    'https://randomuser.me/api/portraits/women/89.jpg'
),

-- Developer 9: Senior Rust Engineer
(
    9, 
    '0430a9b0f2a0bad383b1b3a1989571b90f7486a86629e040c603f6f9ecec857505fd2b1279ccce579dbe59cc88d8d49b7543bd62051b1417cafa6bb2e4fd011d30', 
    'Fred Saloni',
    'fredsaloni',
    'Senior Rust Engineer with focus on building complex system. I build solutions that scale'
    'https://randomuser.me/api/portraits/men/98.jpg'
);