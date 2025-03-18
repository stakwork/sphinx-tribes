-- Insert dummy data for bounties leaderboard
INSERT INTO bounty (
    owner_id, paid, show, completed, type, award, assigned_hours, 
    bounty_expires, commitment_fee, price, title, tribe, assignee, 
    ticket_url, workspace_uuid, feature_uuid, description, wanted_type, 
    deliverables, github_description, one_sentence_summary, 
    estimated_session_length, estimated_completion_date, created, 
    updated, assigned_date, completion_date, mark_as_paid_date, 
    paid_date, coding_languages, phase_uuid, phase_priority,
    payment_pending, payment_failed
) VALUES
-- Developer 1 Bounties
(
    'owner123', true, true, true, 'feature', 'cash', 10, 
    '2023-12-15', 5000, 150000, 'Implement Authentication System', 'backend-tribe', 'developer_1', 
    'https://github.com/org/repo/issues/1', 'workspace-uuid-123', 'feature-uuid-123', 
    'Implement a secure authentication system with JWT tokens', 'backend', 
    'Working JWT authentication with refresh tokens', true, 'Add secure authentication to API', 
    '3 days', '2023-11-25', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '28 days'))::bigint,
    NOW() - INTERVAL '26 days', 
    NOW() - INTERVAL '25 days', 
    NOW() - INTERVAL '22 days', 
    NOW() - INTERVAL '21 days', 
    NOW() - INTERVAL '20 days', 
    ARRAY['Go', 'JavaScript', 'SQL'], 'phase-123', 1, false, false
),
(
    'owner123', true, true, true, 'bug', 'cash', 4, 
    '2023-12-20', 2000, 75000, 'Fix Database Connection Leaks', 'backend-tribe', 'developer_1', 
    'https://github.com/org/repo/issues/2', 'workspace-uuid-123', 'feature-uuid-124', 
    'Fix connection pool issues causing memory leaks', 'backend', 
    'Optimized connection handling with proper resource cleanup', true, 'Solve database connection leaks', 
    '1 day', '2023-11-22', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '26 days'))::bigint,
    NOW() - INTERVAL '24 days', 
    NOW() - INTERVAL '23 days', 
    NOW() - INTERVAL '20 days', 
    NOW() - INTERVAL '19 days', 
    NOW() - INTERVAL '18 days', 
    ARRAY['Go', 'SQL'], 'phase-124', 2, false, false
),
(
    'owner456', true, true, true, 'enhancement', 'cash', 6, 
    '2023-12-25', 3000, 95000, 'Performance Optimization for API', 'backend-tribe', 'developer_1', 
    'https://github.com/org/repo/issues/3', 'workspace-uuid-123', 'feature-uuid-125', 
    'Optimize API response time by implementing caching', 'backend', 
    'Redis caching layer implementation with cache invalidation', true, 'Make API faster with caching', 
    '2 days', '2023-11-28', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '24 days'))::bigint,
    NOW() - INTERVAL '22 days', 
    NOW() - INTERVAL '21 days', 
    NOW() - INTERVAL '18 days', 
    NOW() - INTERVAL '17 days', 
    NOW() - INTERVAL '16 days', 
    ARRAY['Go', 'Redis'], 'phase-125', 1, false, false
),

-- Developer 2 Bounties
(
    'owner789', true, true, true, 'feature', 'cash', 12, 
    '2023-12-18', 6000, 180000, 'Build Real-time Chat Module', 'frontend-tribe', 'developer_2', 
    'https://github.com/org/repo/issues/4', 'workspace-uuid-124', 'feature-uuid-126', 
    'Create a WebSocket-based real-time chat system', 'fullstack', 
    'Working chat system with message persistence', true, 'Add real-time chat functionality', 
    '4 days', '2023-11-26', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '22 days'))::bigint,
    NOW() - INTERVAL '20 days', 
    NOW() - INTERVAL '19 days', 
    NOW() - INTERVAL '15 days', 
    NOW() - INTERVAL '14 days', 
    NOW() - INTERVAL '13 days', 
    ARRAY['JavaScript', 'React', 'Go'], 'phase-126', 1, false, false
),
(
    'owner123', true, true, true, 'bug', 'cash', 3, 
    '2023-12-22', 1500, 65000, 'Fix Cross-browser Compatibility Issues', 'frontend-tribe', 'developer_2', 
    'https://github.com/org/repo/issues/5', 'workspace-uuid-124', 'feature-uuid-127', 
    'Resolve rendering issues in Safari and Firefox', 'frontend', 
    'Working UI across all major browsers', true, 'Fix browser compatibility bugs', 
    '1 day', '2023-11-24', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '20 days'))::bigint,
    NOW() - INTERVAL '18 days', 
    NOW() - INTERVAL '17 days', 
    NOW() - INTERVAL '15 days', 
    NOW() - INTERVAL '14 days', 
    NOW() - INTERVAL '13 days', 
    ARRAY['JavaScript', 'CSS', 'HTML'], 'phase-127', 3, false, false
),

-- Developer 3 Bounties
(
    'owner456', true, true, true, 'feature', 'cash', 15, 
    '2023-12-30', 7500, 225000, 'Design and Implement Admin Dashboard', 'design-tribe', 'developer_3', 
    'https://github.com/org/repo/issues/6', 'workspace-uuid-125', 'feature-uuid-128', 
    'Create an intuitive admin dashboard with analytics', 'frontend', 
    'Fully functional admin dashboard with charts and data visualization', true, 'Build comprehensive admin interface', 
    '5 days', '2023-12-01', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '18 days'))::bigint,
    NOW() - INTERVAL '16 days', 
    NOW() - INTERVAL '15 days', 
    NOW() - INTERVAL '11 days', 
    NOW() - INTERVAL '10 days', 
    NOW() - INTERVAL '9 days', 
    ARRAY['React', 'TypeScript', 'CSS'], 'phase-128', 1, false, false
),
(
    'owner789', true, true, true, 'enhancement', 'cash', 8, 
    '2023-12-28', 4000, 120000, 'Implement Dark Mode Support', 'design-tribe', 'developer_3', 
    'https://github.com/org/repo/issues/7', 'workspace-uuid-125', 'feature-uuid-129', 
    'Add dark mode support with theme switching', 'frontend', 
    'Toggleable light/dark themes with persistent user preference', true, 'Add dark mode to the application', 
    '3 days', '2023-11-30', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '16 days'))::bigint,
    NOW() - INTERVAL '14 days', 
    NOW() - INTERVAL '13 days', 
    NOW() - INTERVAL '10 days', 
    NOW() - INTERVAL '9 days', 
    NOW() - INTERVAL '8 days', 
    ARRAY['CSS', 'JavaScript', 'React'], 'phase-129', 2, false, false
),

-- Developer 4 Bounties
(
    'owner123', true, true, true, 'feature', 'cash', 20, 
    '2023-12-25', 10000, 300000, 'Build Machine Learning Recommendation Engine', 'data-tribe', 'developer_4', 
    'https://github.com/org/repo/issues/8', 'workspace-uuid-126', 'feature-uuid-130', 
    'Implement ML-based content recommendation system', 'backend', 
    'Working recommendation API with training pipeline', true, 'Add smart content recommendations', 
    '7 days', '2023-12-02', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '14 days'))::bigint,
    NOW() - INTERVAL '12 days', 
    NOW() - INTERVAL '11 days', 
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days', 
    NOW() - INTERVAL '5 days', 
    ARRAY['Python', 'TensorFlow', 'SQL'], 'phase-130', 1, false, false
),

-- Developer 5 Bounties
(
    'owner456', true, true, true, 'bug', 'cash', 5, 
    '2023-12-15', 2500, 85000, 'Fix Payment Processing Issues', 'payment-tribe', 'developer_5', 
    'https://github.com/org/repo/issues/9', 'workspace-uuid-127', 'feature-uuid-131', 
    'Resolve payment gateway integration issues', 'backend', 
    'Reliable payment processing with proper error handling', true, 'Fix recurring payment failures', 
    '2 days', '2023-11-25', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '13 days'))::bigint,
    NOW() - INTERVAL '11 days', 
    NOW() - INTERVAL '10 days', 
    NOW() - INTERVAL '8 days', 
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days', 
    ARRAY['Go', 'JavaScript', 'SQL'], 'phase-131', 1, false, false
),
(
    'owner789', true, true, true, 'feature', 'cash', 14, 
    '2023-12-20', 7000, 210000, 'Implement Subscription Management System', 'payment-tribe', 'developer_5', 
    'https://github.com/org/repo/issues/10', 'workspace-uuid-127', 'feature-uuid-132', 
    'Build a system to manage user subscriptions and billing', 'fullstack', 
    'Complete subscription management with billing history', true, 'Add subscription management capabilities', 
    '5 days', '2023-11-28', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '12 days'))::bigint,
    NOW() - INTERVAL '10 days', 
    NOW() - INTERVAL '9 days', 
    NOW() - INTERVAL '6 days', 
    NOW() - INTERVAL '5 days', 
    NOW() - INTERVAL '4 days', 
    ARRAY['Go', 'React', 'SQL'], 'phase-132', 2, false, false
),

-- Developer 6 Bounties
(
    'owner123', true, true, true, 'feature', 'cash', 16, 
    '2023-12-22', 8000, 240000, 'Create Interactive Data Visualization Dashboard', 'analytics-tribe', 'developer_6', 
    'https://github.com/org/repo/issues/11', 'workspace-uuid-128', 'feature-uuid-133', 
    'Build interactive charts and data visualizations', 'frontend', 
    'Real-time data visualization dashboard with filtering capabilities', true, 'Add interactive data visualizations', 
    '6 days', '2023-11-30', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '10 days'))::bigint,
    NOW() - INTERVAL '8 days', 
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '4 days', 
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '2 days', 
    ARRAY['JavaScript', 'D3.js', 'React'], 'phase-133', 1, false, false
),

-- Developer 7 Bounties
(
    'owner456', true, true, true, 'enhancement', 'cash', 10, 
    '2023-12-28', 5000, 150000, 'Optimize Database Queries', 'backend-tribe', 'developer_7', 
    'https://github.com/org/repo/issues/12', 'workspace-uuid-129', 'feature-uuid-134', 
    'Improve database query performance for analytics dashboard', 'backend', 
    'Optimized queries with proper indexing strategy', true, 'Speed up slow database queries', 
    '4 days', '2023-12-03', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '9 days'))::bigint,
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days', 
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '1 day', 
    ARRAY['SQL', 'Go', 'PostgreSQL'], 'phase-134', 1, false, false
),
(
    'owner789', true, true, true, 'feature', 'cash', 12, 
    '2023-12-30', 6000, 180000, 'Implement API Rate Limiting', 'security-tribe', 'developer_7', 
    'https://github.com/org/repo/issues/13', 'workspace-uuid-129', 'feature-uuid-135', 
    'Add rate limiting to protect API endpoints from abuse', 'backend', 
    'Configurable rate limiting with proper response headers', true, 'Protect API from excessive usage', 
    '4 days', '2023-12-05', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '8 days'))::bigint,
    NOW() - INTERVAL '6 days', 
    NOW() - INTERVAL '5 days', 
    NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '12 hours', 
    ARRAY['Go', 'Redis'], 'phase-135', 2, false, false
),

-- Developer 8 Bounties
(
    'owner123', true, true, true, 'bug', 'cash', 7, 
    '2023-12-20', 3500, 105000, 'Fix Mobile Responsiveness Issues', 'frontend-tribe', 'developer_8', 
    'https://github.com/org/repo/issues/14', 'workspace-uuid-130', 'feature-uuid-136', 
    'Resolve UI layout problems on mobile devices', 'frontend', 
    'Fully responsive UI across all screen sizes', true, 'Fix mobile UI bugs', 
    '3 days', '2023-11-26', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '7 days'))::bigint,
    NOW() - INTERVAL '5 days', 
    NOW() - INTERVAL '4 days', 
    NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '12 hours', 
    ARRAY['CSS', 'HTML', 'JavaScript'], 'phase-136', 3, false, false
),
(
    'owner456', true, true, true, 'feature', 'cash', 15, 
    '2023-12-25', 7500, 225000, 'Implement Multi-language Support', 'frontend-tribe', 'developer_8', 
    'https://github.com/org/repo/issues/15', 'workspace-uuid-130', 'feature-uuid-137', 
    'Add internationalization and localization support', 'fullstack', 
    'UI and content translation for multiple languages', true, 'Make application support multiple languages', 
    '5 days', '2023-12-01', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '6 days'))::bigint,
    NOW() - INTERVAL '4 days', 
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '12 hours', 
    NOW() - INTERVAL '6 hours', 
    ARRAY['JavaScript', 'React', 'i18n'], 'phase-137', 1, false, false
),
(
    'owner789', true, true, true, 'enhancement', 'cash', 8, 
    '2023-12-28', 4000, 120000, 'Improve Accessibility Compliance', 'frontend-tribe', 'developer_8', 
    'https://github.com/org/repo/issues/16', 'workspace-uuid-130', 'feature-uuid-138', 
    'Enhance application to meet WCAG 2.1 AA standards', 'frontend', 
    'Fully accessible UI with screen reader support', true, 'Make app accessible to all users', 
    '3 days', '2023-12-03', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '5 days'))::bigint,
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '12 hours', 
    NOW() - INTERVAL '6 hours', 
    ARRAY['HTML', 'CSS', 'JavaScript'], 'phase-138', 2, false, false
);

-- Add more bounties for each developer to create variety

-- Developer 1 additional bounties
INSERT INTO bounty (
    owner_id, paid, show, completed, type, award, assigned_hours, 
    bounty_expires, commitment_fee, price, title, tribe, assignee, 
    ticket_url, workspace_uuid, feature_uuid, description, wanted_type, 
    deliverables, github_description, one_sentence_summary, 
    estimated_session_length, estimated_completion_date, created, 
    updated, assigned_date, completion_date, mark_as_paid_date, 
    paid_date, coding_languages, phase_uuid, phase_priority,
    payment_pending, payment_failed
) VALUES
(
    'owner789', true, true, true, 'feature', 'cash', 18, 
    '2023-12-15', 9000, 270000, 'Build Microservice Architecture', 'architecture-tribe', 'developer_1', 
    'https://github.com/org/repo/issues/17', 'workspace-uuid-131', 'feature-uuid-139', 
    'Refactor monolith into microservices architecture', 'backend', 
    'Working microservices with service discovery and API gateway', true, 'Modernize application architecture', 
    '6 days', '2023-12-05', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '25 days'))::bigint,
    NOW() - INTERVAL '23 days', 
    NOW() - INTERVAL '22 days', 
    NOW() - INTERVAL '17 days', 
    NOW() - INTERVAL '16 days', 
    NOW() - INTERVAL '15 days', 
    ARRAY['Go', 'Docker', 'Kubernetes'], 'phase-139', 1, false, false
),
(
    'owner123', true, true, true, 'enhancement', 'cash', 9, 
    '2023-12-20', 4500, 135000, 'Implement OAuth2 Authentication', 'security-tribe', 'developer_1', 
    'https://github.com/org/repo/issues/18', 'workspace-uuid-131', 'feature-uuid-140', 
    'Add OAuth2 support for third-party authentication', 'backend', 
    'Working OAuth2 integration with multiple providers', true, 'Add social login functionality', 
    '3 days', '2023-11-28', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '23 days'))::bigint,
    NOW() - INTERVAL '21 days', 
    NOW() - INTERVAL '20 days', 
    NOW() - INTERVAL '16 days', 
    NOW() - INTERVAL '15 days', 
    NOW() - INTERVAL '14 days', 
    ARRAY['Go', 'JavaScript'], 'phase-140', 2, false, false
);

-- Developer 2 additional bounties  
INSERT INTO bounty (
    owner_id, paid, show, completed, type, award, assigned_hours, 
    bounty_expires, commitment_fee, price, title, tribe, assignee, 
    ticket_url, workspace_uuid, feature_uuid, description, wanted_type, 
    deliverables, github_description, one_sentence_summary, 
    estimated_session_length, estimated_completion_date, created, 
    updated, assigned_date, completion_date, mark_as_paid_date, 
    paid_date, coding_languages, phase_uuid, phase_priority,
    payment_pending, payment_failed
) VALUES
(
    'owner456', true, true, true, 'feature', 'cash', 14, 
    '2023-12-22', 7000, 210000, 'Create Interactive Form Builder', 'frontend-tribe', 'developer_2', 
    'https://github.com/org/repo/issues/19', 'workspace-uuid-132', 'feature-uuid-141', 
    'Build a drag-and-drop form builder with validation', 'frontend', 
    'Working form builder with custom validation rules', true, 'Add form builder functionality', 
    '5 days', '2023-11-30', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '19 days'))::bigint,
    NOW() - INTERVAL '17 days', 
    NOW() - INTERVAL '16 days', 
    NOW() - INTERVAL '12 days', 
    NOW() - INTERVAL '11 days', 
    NOW() - INTERVAL '10 days', 
    ARRAY['React', 'TypeScript', 'CSS'], 'phase-141', 1, false, false
),
(
    'owner789', true, true, true, 'enhancement', 'cash', 7, 
    '2023-12-25', 3500, 105000, 'Implement Lazy Loading for Images', 'performance-tribe', 'developer_2', 
    'https://github.com/org/repo/issues/20', 'workspace-uuid-132', 'feature-uuid-142', 
    'Add lazy loading to improve initial page load time', 'frontend', 
    'Optimized image loading with placeholder support', true, 'Improve image loading performance', 
    '2 days', '2023-12-01', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '17 days'))::bigint,
    NOW() - INTERVAL '15 days', 
    NOW() - INTERVAL '14 days', 
    NOW() - INTERVAL '11 days', 
    NOW() - INTERVAL '10 days', 
    NOW() - INTERVAL '9 days', 
    ARRAY['JavaScript', 'HTML', 'CSS'], 'phase-142', 2, false, false
);

-- More random bounties spread across all developers
INSERT INTO bounty (
    owner_id, paid, show, completed, type, award, assigned_hours, 
    bounty_expires, commitment_fee, price, title, tribe, assignee, 
    ticket_url, workspace_uuid, feature_uuid, description, wanted_type, 
    deliverables, github_description, one_sentence_summary, 
    estimated_session_length, estimated_completion_date, created, 
    updated, assigned_date, completion_date, mark_as_paid_date, 
    paid_date, coding_languages, phase_uuid, phase_priority,
    payment_pending, payment_failed
) VALUES
(
    'owner123', true, true, true, 'feature', 'cash', 25, 
    '2023-12-30', 12500, 375000, 'Build CI/CD Pipeline', 'devops-tribe', 'developer_3', 
    'https://github.com/org/repo/issues/21', 'workspace-uuid-133', 'feature-uuid-143', 
    'Implement automated CI/CD pipeline for deployment', 'devops', 
    'Working pipeline with testing, building and deployment', true, 'Automate the deployment process', 
    '8 days', '2023-12-10', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '15 days'))::bigint,
    NOW() - INTERVAL '13 days', 
    NOW() - INTERVAL '12 days', 
    NOW() - INTERVAL '8 days', 
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days', 
    ARRAY['Docker', 'Jenkins', 'YAML'], 'phase-143', 1, false, false
),
(
    'owner456', true, true, true, 'bug', 'cash', 6, 
    '2023-12-15', 3000, 90000, 'Fix Memory Leaks in Frontend', 'performance-tribe', 'developer_4', 
    'https://github.com/org/repo/issues/22', 'workspace-uuid-134', 'feature-uuid-144', 
    'Resolve memory leaks in React components', 'frontend', 
    'Optimized components with no memory leaks', true, 'Fix application memory issues', 
    '2 days', '2023-11-25', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '13 days'))::bigint,
    NOW() - INTERVAL '11 days', 
    NOW() - INTERVAL '10 days', 
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days', 
    NOW() - INTERVAL '5 days', 
    ARRAY['JavaScript', 'React'], 'phase-144', 2, false, false
),
(
    'owner789', true, true, true, 'feature', 'cash', 20, 
    '2023-12-28', 10000, 300000, 'Implement Notifications System', 'backend-tribe', 'developer_5', 
    'https://github.com/org/repo/issues/23', 'workspace-uuid-135', 'feature-uuid-145', 
    'Build real-time notification system with multiple channels', 'fullstack', 
    'Working notifications via WebSocket, email, and mobile push', true, 'Add comprehensive notification capabilities', 
    '7 days', '2023-12-05', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '11 days'))::bigint,
    NOW() - INTERVAL '9 days', 
    NOW() - INTERVAL '8 days', 
    NOW() - INTERVAL '4 days', 
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '2 days', 
    ARRAY['Go', 'React', 'WebSocket'], 'phase-145', 1, false, false
),
(
    'owner123', true, true, true, 'enhancement', 'cash', 12, 
    '2023-12-20', 6000, 180000, 'Implement Content Delivery Network', 'infrastructure-tribe', 'developer_6', 
    'https://github.com/org/repo/issues/24', 'workspace-uuid-136', 'feature-uuid-146', 
    'Set up CDN for static assets to improve load times', 'devops', 
    'Working CDN integration with proper cache configuration', true, 'Speed up content delivery globally', 
    '4 days', '2023-11-30', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '9 days'))::bigint,
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days', 
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '1 day', 
    ARRAY['AWS', 'CloudFront', 'Terraform'], 'phase-146', 2, false, false
),
(
    'owner456', true, true, true, 'feature', 'cash', 18, 
    '2023-12-25', 9000, 270000, 'Build Analytics Dashboard', 'data-tribe', 'developer_7', 
    'https://github.com/org/repo/issues/25', 'workspace-uuid-137', 'feature-uuid-147', 
    'Create a dashboard to visualize user behavior analytics', 'fullstack', 
    'Interactive analytics dashboard with filtering capabilities', true, 'Add analytics visualization capabilities', 
    '6 days', '2023-12-03', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '7 days'))::bigint,
    NOW() - INTERVAL '5 days', 
    NOW() - INTERVAL '4 days', 
    NOW() - INTERVAL '2 days', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '12 hours', 
    ARRAY['React', 'D3.js', 'Go'], 'phase-147', 1, false, false
),
(
    'owner789', true, true, true, 'bug', 'cash', 8, 
    '2023-12-18', 4000, 120000, 'Fix SEO Issues', 'marketing-tribe', 'developer_8', 
    'https://github.com/org/repo/issues/26', 'workspace-uuid-138', 'feature-uuid-148', 
    'Resolve SEO problems and improve search engine ranking', 'frontend', 
    'Improved SEO with proper meta tags and structured data', true, 'Boost search engine visibility', 
    '3 days', '2023-11-28', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '6 days'))::bigint,
    NOW() - INTERVAL '4 days', 
    NOW() - INTERVAL '3 days', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '12 hours', 
    NOW() - INTERVAL '6 hours', 
    ARRAY['HTML', 'JavaScript', 'SEO'], 'phase-148', 3, false, false
),
(
    'owner123', true, true, true, 'feature', 'cash', 25, 
    '2023-12-30', 12500, 375000, 'Setup a Rust Rocket Web Server', 'rust-server', '0430a9b0f2a0bad383b1b3a1989571b90f7486a86629e040c603f6f9ecec857505fd2b1279ccce579dbe59cc88d8d49b7543bd62051b1417cafa6bb2e4fd011d30', 
    'https://github.com/org/repo/issues/21', 'workspace-uuid-133', 'feature-uuid-143', 
    'We need an Http Server the web and mobile app can connect', 'backend', 
    'Setup Unit and Integration test', true, 'Setup the Backend Server', 
    '8 days', '2023-12-10', 
    EXTRACT(EPOCH FROM (NOW() - INTERVAL '15 days'))::bigint,
    NOW() - INTERVAL '13 days', 
    NOW() - INTERVAL '12 days', 
    NOW() - INTERVAL '8 days', 
    NOW() - INTERVAL '7 days', 
    NOW() - INTERVAL '6 days', 
    ARRAY['Rust', 'Rocket', 'YAML'], 'phase-143', 1, false, false
);