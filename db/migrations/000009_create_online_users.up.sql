-- Create online_users table
CREATE TABLE online_users (
    online_user_id  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL,
    socket_id       VARCHAR(50) NOT NULL,
    device_info     VARCHAR(500),
    current_page_url VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255),
    CONSTRAINT fk_online_user_user_id FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT uq_online_users_user_id_socket_id UNIQUE(user_id, socket_id)
);

CREATE INDEX idx_online_users_created_user_id ON online_users(created_user_id);
CREATE INDEX idx_online_users_updated_user_id ON online_users(updated_user_id);
CREATE INDEX idx_online_users_deleted_user_id ON online_users(deleted_user_id);
CREATE INDEX idx_online_users_created_at ON online_users(created_at);
CREATE INDEX idx_online_users_updated_at ON online_users(updated_at);
CREATE INDEX idx_online_users_deleted_at ON online_users(deleted_at);
CREATE INDEX idx_online_users_user_id ON online_users(user_id);
