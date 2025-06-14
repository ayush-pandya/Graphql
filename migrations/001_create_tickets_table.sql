-- Create tickets table
CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    priority VARCHAR(50) NOT NULL DEFAULT 'MEDIUM',
    assignee_id VARCHAR(100),
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on status and priority for faster queries
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority);
CREATE INDEX IF NOT EXISTS idx_tickets_assignee ON tickets(assignee_id);
CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets(created_at);

-- Insert some sample data
INSERT INTO tickets (id, title, description, status, priority, assignee_id, tags) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Fix login bug', 'Users cannot log in with their email', 'OPEN', 'HIGH', 'user-123', ARRAY['bug', 'urgent']),
    ('550e8400-e29b-41d4-a716-446655440002', 'Add dark mode', 'Implement dark mode for better user experience', 'IN_PROGRESS', 'MEDIUM', 'user-456', ARRAY['feature', 'ui']),
    ('550e8400-e29b-41d4-a716-446655440003', 'Performance optimization', 'Optimize database queries for faster response', 'OPEN', 'LOW', 'user-789', ARRAY['performance', 'database'])
ON CONFLICT (id) DO NOTHING; 