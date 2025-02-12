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
    ('1', 'Implement drag and drop functionality', 'Add smooth drag and drop between columns in the Kanban board view', '2026-05-15', 'John', NOW(), NOW(), NULL, 'in_progress'),
    ('2', 'Design user dashboard', 'Create wireframes and mockups for the main dashboard interface', '2026-05-10', 'Laura', NOW(), NOW(), NULL, 'in_review'),
    ('3', 'Add board sharing capabilities', 'Implement feature to share boards with team members with different permission levels', '2026-05-20', 'Steve', NOW(), NOW(), NULL, 'pending'),
    ('4', 'Optimize database queries', 'Improve performance of board loading and card updates', '2026-05-12', 'Alex', NOW(), NOW(), NULL, 'in_progress'),
    ('5', 'Implement user authentication', 'Set up JWT authentication and user session management', '2026-05-08', 'Mary', NOW(), NOW(), '2026-04-20', 'completed'),
    ('6', 'Create board templates', 'Design and implement reusable board templates for common use cases', '2026-05-25', 'John', NOW(), NOW(), NULL, 'pending'),
    ('7', 'Add card labels feature', 'Implement color-coded labels for better task categorization', '2026-05-18', 'Laura', NOW(), NOW(), NULL, 'in_progress'),
    ('8', 'Mobile responsive design', 'Ensure proper display and functionality on mobile devices', '2026-05-30', 'Steve', NOW(), NOW(), NULL, 'pending'),
    ('9', 'Implement activity log', 'Add tracking for all board and card changes with user attribution', '2026-05-22', 'Alex', NOW(), NOW(), NULL, 'in_review');