-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE ON UPDATE CASCADE,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER posts_update_at_trigger
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(updated_at);

-- +goose Down

DROP TRIGGER posts_update_at_trigger ON posts;
DROP TABLE posts;
