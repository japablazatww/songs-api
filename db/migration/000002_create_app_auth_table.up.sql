-- Up migration
CREATE TABLE app_auth (
    id BIGSERIAL PRIMARY KEY,
    app_name VARCHAR(255) NOT NULL,
    client_id VARCHAR(32) NOT NULL UNIQUE,
    client_secret VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
