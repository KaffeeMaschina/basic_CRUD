-- +goose Up
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementBegin
SELECT 'up SQL query';

-- +goose StatementEnd

-- +goose Down
    DROP TABLE tasks;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
