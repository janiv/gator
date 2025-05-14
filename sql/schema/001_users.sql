-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;