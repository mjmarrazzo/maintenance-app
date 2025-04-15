-- Create enums for statuses and priorities
CREATE TYPE task_priority AS ENUM ('Low', 'Medium', 'High', 'Urgent');
CREATE TYPE task_status AS ENUM ('New', 'In Progress', 'Completed', 'On Hold');
CREATE TYPE user_role AS ENUM ('Member', 'Maintenance', 'Administrator');
CREATE TYPE recurrence_type AS ENUM ('Daily', 'Weekly', 'Monthly', 'Yearly', 'Custom');
CREATE TYPE recurrence_unit AS ENUM ('Days', 'Weeks', 'Months', 'Years');

-- Create Categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

-- Create Locations table with hierarchical structure
CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_location_id INTEGER REFERENCES locations(id) ON DELETE CASCADE
);

-- Create Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'Member',
    phone VARCHAR(20),
    notification_preferences JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create Tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    location_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,
    priority task_priority NOT NULL DEFAULT 'Medium',
    status task_status NOT NULL DEFAULT 'New',
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_to INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    estimated_completion_date TIMESTAMP WITH TIME ZONE,
    cost DECIMAL(10, 2),
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    recurrence_type recurrence_type NULL,
    recurrence_interval INTEGER,
    recurrence_unit recurrence_unit NULL,
    parent_task_id INTEGER REFERENCES tasks(id) ON DELETE SET NULL,
    next_occurrence TIMESTAMP WITH TIME ZONE
);

-- Create Comments table
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create CommentReactions table
CREATE TABLE IF NOT EXISTS comment_reactions (
    id SERIAL PRIMARY KEY,
    comment_id INTEGER NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji_code VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(comment_id, user_id, emoji_code)
);

-- Create Attachments table
CREATE TABLE IF NOT EXISTS attachments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    object_key VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size_bytes INTEGER NOT NULL,
    uploaded_by INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    upload_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create TaskHistory table for audit trail
CREATE TABLE IF NOT EXISTS task_history (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    changed_by INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    changed_field VARCHAR(50) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance optimization
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_category_id ON tasks(category_id);
CREATE INDEX idx_tasks_location_id ON tasks(location_id);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);
CREATE INDEX idx_tasks_is_recurring ON tasks(is_recurring);
CREATE INDEX idx_comments_task_id ON comments(task_id);
CREATE INDEX idx_attachments_task_id ON attachments(task_id);
CREATE INDEX idx_task_history_task_id ON task_history(task_id);

-- Create trigger to update updated_at timestamp on tasks
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tasks_modtime
BEFORE UPDATE ON tasks
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

-- Create trigger to log task changes to task_history
CREATE OR REPLACE FUNCTION log_task_changes()
RETURNS TRIGGER AS $$
DECLARE
    changed_field TEXT;
    old_value TEXT;
    new_value TEXT;
BEGIN
    IF TG_OP = 'UPDATE' THEN
        -- Check each field for changes and log them
        IF OLD.title IS DISTINCT FROM NEW.title THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'title', OLD.title, NEW.title);
        END IF;

        IF OLD.description IS DISTINCT FROM NEW.description THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'description', OLD.description, NEW.description);
        END IF;

        IF OLD.category_id IS DISTINCT FROM NEW.category_id THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'category_id', OLD.category_id::TEXT, NEW.category_id::TEXT);
        END IF;

        IF OLD.location_id IS DISTINCT FROM NEW.location_id THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'location_id', OLD.location_id::TEXT, NEW.location_id::TEXT);
        END IF;

        IF OLD.priority IS DISTINCT FROM NEW.priority THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'priority', OLD.priority::TEXT, NEW.priority::TEXT);
        END IF;

        IF OLD.status IS DISTINCT FROM NEW.status THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'status', OLD.status::TEXT, NEW.status::TEXT);
        END IF;

        IF OLD.assigned_to IS DISTINCT FROM NEW.assigned_to THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'assigned_to', OLD.assigned_to::TEXT, NEW.assigned_to::TEXT);
        END IF;

        IF OLD.estimated_completion_date IS DISTINCT FROM NEW.estimated_completion_date THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'estimated_completion_date', OLD.estimated_completion_date::TEXT, NEW.estimated_completion_date::TEXT);
        END IF;

        IF OLD.cost IS DISTINCT FROM NEW.cost THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'cost', OLD.cost::TEXT, NEW.cost::TEXT);
        END IF;

        IF OLD.is_recurring IS DISTINCT FROM NEW.is_recurring THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'is_recurring', OLD.is_recurring::TEXT, NEW.is_recurring::TEXT);
        END IF;

        IF OLD.recurrence_pattern IS DISTINCT FROM NEW.recurrence_pattern THEN
            INSERT INTO task_history (task_id, changed_by, changed_field, old_value, new_value)
            VALUES (NEW.id, COALESCE(NEW.assigned_to, NEW.created_by), 'recurrence_pattern', OLD.recurrence_pattern::TEXT, NEW.recurrence_pattern::TEXT);
        END IF;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER task_history_trigger
AFTER UPDATE ON tasks
FOR EACH ROW EXECUTE FUNCTION log_task_changes();

-- Insert default categories
INSERT INTO categories (name, description)
VALUES
('Electrical', 'Electrical system maintenance and repairs'),
('Plumbing', 'Plumbing system maintenance and repairs'),
('Cleaning', 'Regular and deep cleaning tasks'),
('Structural', 'Building structure maintenance'),
('HVAC', 'Heating, ventilation, and air conditioning'),
('Grounds', 'Outdoor areas and landscaping'),
('Technology', 'IT systems and audio/visual equipment'),
('Furniture', 'Furniture repair and maintenance')
ON CONFLICT DO NOTHING;

-- Insert a default admin user with password 'admin' (this is just for initial setup, should be changed)
-- Note: In production, use a proper password hashing mechanism
INSERT INTO users (name, email, password_hash, role)
VALUES ('Admin', 'admin@churchmaintenance.org', '$2a$10$9tWMVVIQAj8SHfDFaRpp4e5AsvSo5dP4OmYjt3HhKkIeJhAYMrE1u', 'Administrator')
ON CONFLICT DO NOTHING;

-- Function to calculate next occurrence
CREATE OR REPLACE FUNCTION calculate_next_occurrence()
RETURNS TRIGGER AS $$
BEGIN
    -- Only calculate for recurring tasks
    IF NEW.is_recurring = TRUE THEN
        -- If next_occurrence is null or already passed, calculate new one
        IF NEW.next_occurrence IS NULL OR NEW.next_occurrence < NOW() THEN
            -- Calculate based on recurrence type
            CASE NEW.recurrence_type
                WHEN 'Daily' THEN
                    NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                        (NEW.recurrence_interval || ' days')::INTERVAL;

                WHEN 'Weekly' THEN
                    NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                        (NEW.recurrence_interval * 7 || ' days')::INTERVAL;

                WHEN 'Monthly' THEN
                    NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                        (NEW.recurrence_interval || ' months')::INTERVAL;

                WHEN 'Yearly' THEN
                    NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                        (NEW.recurrence_interval || ' years')::INTERVAL;

                WHEN 'Custom' THEN
                    -- Handle custom recurrence using recurrence_unit
                    CASE NEW.recurrence_unit
                        WHEN 'Days' THEN
                            NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                                (NEW.recurrence_interval || ' days')::INTERVAL;

                        WHEN 'Weeks' THEN
                            NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                                (NEW.recurrence_interval * 7 || ' days')::INTERVAL;

                        WHEN 'Months' THEN
                            NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                                (NEW.recurrence_interval || ' months')::INTERVAL;

                        WHEN 'Years' THEN
                            NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) +
                                (NEW.recurrence_interval || ' years')::INTERVAL;

                        ELSE
                            -- Default fallback
                            NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) + '1 day'::INTERVAL;
                    END CASE;

                ELSE
                    -- Default fallback
                    NEW.next_occurrence := COALESCE(NEW.next_occurrence, NOW()) + '1 day'::INTERVAL;
            END CASE;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to run before insert or update
CREATE TRIGGER update_next_occurrence
BEFORE INSERT OR UPDATE ON tasks
FOR EACH ROW EXECUTE FUNCTION calculate_next_occurrence();

-- Create a function to update all outdated occurrences
-- This can be called by a scheduled job
CREATE OR REPLACE FUNCTION update_outdated_occurrences()
RETURNS INTEGER AS $$
DECLARE
    updated_count INTEGER := 0;
BEGIN
    -- Update all tasks with outdated next_occurrence dates
    WITH updated_rows AS (
        UPDATE tasks
        SET next_occurrence = next_occurrence -- This will trigger the calculate_next_occurrence function
        WHERE is_recurring = TRUE
        AND next_occurrence < NOW()
        RETURNING id
    )
    SELECT COUNT(*) INTO updated_count FROM updated_rows;

    RETURN updated_count;
END;
$$ LANGUAGE plpgsql;