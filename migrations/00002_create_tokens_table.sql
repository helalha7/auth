-- +goose Up
CREATE TABLE tokens (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,

    hash BINARY(32) NOT NULL UNIQUE,
    expiry TIMESTAMP NOT NULL,
    scope VARCHAR(200) NOT NULL,

    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE tokens;
