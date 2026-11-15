-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id       SERIAL PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS reviews (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id),
    text       TEXT NOT NULL,
    star_count INT NOT NULL CHECK (star_count BETWEEN 1 AND 5)
);

-- +goose Down
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS users;
