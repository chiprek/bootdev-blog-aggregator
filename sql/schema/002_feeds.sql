-- +goose Up

CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TRIGGER feeds_update_at_trigger
    BEFORE UPDATE ON feeds
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(updated_at);

-- +goose Down
DROP TRIGGER feeds_update_at_trigger ON feeds;
DROP TABLE feeds;
