-- Create user_roles table
CREATE TABLE user_roles (
    user_role_id    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL,
    role_id         UUID        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255),
    CONSTRAINT fk_user_role_user_id FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role_role_id FOREIGN KEY(role_id) REFERENCES roles(role_id) ON DELETE CASCADE
);

CREATE INDEX idx_user_roles_created_user_id ON user_roles(created_user_id);
CREATE INDEX idx_user_roles_updated_user_id ON user_roles(updated_user_id);
CREATE INDEX idx_user_roles_deleted_user_id ON user_roles(deleted_user_id);
CREATE INDEX idx_user_roles_created_at ON user_roles(created_at);
CREATE INDEX idx_user_roles_updated_at ON user_roles(updated_at);
CREATE INDEX idx_user_roles_deleted_at ON user_roles(deleted_at);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
