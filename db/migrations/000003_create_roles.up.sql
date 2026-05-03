-- Create roles table
CREATE TABLE roles (
    role_id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    description     VARCHAR(255),
    module_id       UUID        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255),
    CONSTRAINT fk_role_module_id FOREIGN KEY(module_id) REFERENCES modules(module_id) ON DELETE CASCADE
);

CREATE INDEX idx_roles_created_user_id ON roles(created_user_id);
CREATE INDEX idx_roles_updated_user_id ON roles(updated_user_id);
CREATE INDEX idx_roles_deleted_user_id ON roles(deleted_user_id);
CREATE INDEX idx_roles_created_at ON roles(created_at);
CREATE INDEX idx_roles_updated_at ON roles(updated_at);
CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);
CREATE INDEX idx_roles_module_id ON roles(module_id);
