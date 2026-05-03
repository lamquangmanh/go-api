-- Create resources table
CREATE TABLE resources (
    resource_id     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    module_id       UUID        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255),
    CONSTRAINT fk_resource_module_id FOREIGN KEY(module_id) REFERENCES modules(module_id) ON DELETE CASCADE
);

CREATE INDEX idx_resources_created_user_id ON resources(created_user_id);
CREATE INDEX idx_resources_updated_user_id ON resources(updated_user_id);
CREATE INDEX idx_resources_deleted_user_id ON resources(deleted_user_id);
CREATE INDEX idx_resources_created_at ON resources(created_at);
CREATE INDEX idx_resources_updated_at ON resources(updated_at);
CREATE INDEX idx_resources_deleted_at ON resources(deleted_at);
CREATE INDEX idx_resources_module_id ON resources(module_id);
