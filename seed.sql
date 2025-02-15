INSERT INTO tasks (
    id, 
    title, 
    description, 
    due_date, 
    assignee_name, 
    created_at, 
    updated_at, 
    completed_at, 
    status
) VALUES 
    ('1', 'Implement Biometric Authentication', 'Add fingerprint and Face ID support for quick login across iOS and Android platforms', '2026-05-15', 'John', NOW(), NOW(), NULL, 'in_progress'),
    ('2', 'Design Transaction Dashboard', 'Create intuitive visualization for spending patterns, investment portfolio, and bill payment reminders', '2026-05-10', 'Laura', NOW(), NOW(), NULL, 'in_review'),
    ('3', 'Develop Peer Transfer System', 'Build secure P2P money transfer feature with contact integration and transaction limits', '2026-05-20', 'Steve', NOW(), NOW(), NULL, 'pending'),
    ('4', 'Optimize Payment Processing', 'Improve transaction speed and reduce latency in payment confirmation system', '2026-05-12', 'Alex', NOW(), NOW(), NULL, 'in_progress'),
    ('5', 'Implement KYC Verification', 'Set up automated document verification and identity validation workflow', '2026-05-08', 'Mary', NOW(), NOW(), '2026-04-20', 'completed'),
    ('6', 'Create Investment Templates', 'Design preset investment portfolios for different risk profiles and goals', '2026-05-25', 'John', NOW(), NOW(), NULL, 'pending'),
    ('7', 'Add Smart Bill Categories', 'Implement ML-based transaction categorization and spending insights', '2026-05-18', 'Laura', NOW(), NOW(), NULL, 'in_progress'),
    ('8', 'Develop Offline Transaction Mode', 'Enable basic functionality and transaction queueing without internet connection', '2026-05-30', 'Steve', NOW(), NOW(), NULL, 'pending'),
    ('9', 'Implement Transaction History', 'Add detailed transaction logging with advanced search and export capabilities', '2026-05-22', 'Alex', NOW(), NOW(), NULL, 'in_review');