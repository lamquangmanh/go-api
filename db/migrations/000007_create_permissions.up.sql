-- Create permissions table
CREATE TABLE permissions (
    permission_id   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id         UUID        NOT NULL,
    resource_id     UUID        NOT NULL,
    action_id       UUID        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255),
    CONSTRAINT fk_permission_role_id FOREIGN KEY(role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    CONSTRAINT fk_permission_resource_id FOREIGN KEY(resource_id) REFERENCES resources(resource_id) ON DELETE CASCADE,
    CONSTRAINT fk_permission_action_id FOREIGN KEY(action_id) REFERENCES actions(action_id) ON DELETE CASCADE
);

CREATE INDEX idx_permissions_created_user_id ON permissions(created_user_id);
CREATE INDEX idx_permissions_updated_user_id ON permissions(updated_user_id);
CREATE INDEX idx_permissions_deleted_user_id ON permissions(deleted_user_id);
CREATE INDEX idx_permissions_created_at ON permissions(created_at);
CREATE INDEX idx_permissions_updated_at ON permissions(updated_at);
CREATE INDEX idx_permissions_deleted_at ON permissions(deleted_at);
CREATE INDEX idx_permissions_role_id ON permissions(role_id);
CREATE INDEX idx_permissions_resource_id ON permissions(resource_id);
CREATE INDEX idx_permissions_action_id ON permissions(action_id);
