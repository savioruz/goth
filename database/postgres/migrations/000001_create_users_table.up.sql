BEGIN;

CREATE TABLE IF NOT EXISTS `users` (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_users_email ON user(email);
CREATE INDEX idx_users_deleted_at ON user USING btree (deleted_at ASC NULLS LAST);

COMMIT;