-- +goose Up

CREATE EXTENSION IF NOT EXISTS moddatetime;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE
);

CREATE TRIGGER users_updated_at_trigger
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE PROCEDURE moddatetime(updated_at);

-- +goose Down
DROP TRIGGER users_updated_at_trigger ON users;
DROP TABLE users;
