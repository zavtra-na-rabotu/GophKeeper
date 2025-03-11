CREATE TABLE secrets
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGSERIAL NOT NULL REFERENCES users ON DELETE CASCADE,
    title      TEXT,
    type       TEXT CHECK (type IN ('credential', 'text', 'binary', 'card')),
    content    BYTEA     NOT NULL,
    metadata   TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);