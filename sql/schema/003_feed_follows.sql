-- +goose Up
CREATE TABLE feed_follows(
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (feed_id, user_id)
);

CREATE TRIGGER feed_follows_updated_at_trigger
    BEFORE UPDATE ON feed_follows
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(updated_at);

-- +goose Down
DROP TRIGGER feed_follows_updated_at_trigger ON feed_follows;
DROP TABLE feed_follows;
