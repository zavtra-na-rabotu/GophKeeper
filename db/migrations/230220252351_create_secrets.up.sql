CREATE TABLE secrets
(
    id         SERIAL PRIMARY KEY,
    user_id    SERIAL NOT NULL REFERENCES users ON DELETE CASCADE,
    title      TEXT,
    type       TEXT CHECK (type IN ('credential', 'text', 'binary', 'card')),
    content    BYTEA  NOT NULL,
    metadata   TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trigger_update_timestamp
    BEFORE UPDATE
    ON secrets
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();