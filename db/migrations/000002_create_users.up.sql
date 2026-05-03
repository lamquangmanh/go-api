-- Create enum types for users
CREATE TYPE user_status AS ENUM ('ACTIVE', 'DEACTIVATED', 'DELETED');

-- Create users table
CREATE TABLE users (
    user_id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_name       VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL,
    password        VARCHAR(255) NOT NULL,
    phone           VARCHAR(20),
    avatar          VARCHAR(1000),
    status          user_status NOT NULL DEFAULT 'ACTIVE',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255)
);

-- Create unique constraint on email and deleted_at (for soft delete)
CREATE UNIQUE INDEX uq_users_email_deleted_at ON users(email, deleted_at)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_users_created_user_id ON users(created_user_id);
CREATE INDEX idx_users_updated_user_id ON users(updated_user_id);
CREATE INDEX idx_users_deleted_user_id ON users(deleted_user_id);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_updated_at ON users(updated_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
