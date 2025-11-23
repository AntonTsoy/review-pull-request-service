CREATE TABLE IF NOT EXISTS teams (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    team_id INTEGER NOT NULL REFERENCES teams (id) ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_users_team_id ON users (team_id);

CREATE TABLE IF NOT EXISTS pull_requests (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author_id VARCHAR(36) NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    status status_enum NOT NULL,
    merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reviewers (
    pull_request_id VARCHAR(50) REFERENCES pull_requests (id) ON DELETE CASCADE,
    reviewer_id VARCHAR(36) REFERENCES users (id) ON DELETE CASCADE,

    PRIMARY KEY (reviewer_id, pull_request_id)
);
