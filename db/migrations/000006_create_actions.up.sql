-- Create enum type for request types
CREATE TYPE request_type AS ENUM ('VIEW', 'HTTP', 'GRAPHQL', 'GRPC', 'WEBSOCKET');

-- Create actions table
CREATE TABLE actions (
    action_id       UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id     UUID        NOT NULL,
    name            VARCHAR(100) NOT NULL,
    description     VARCHAR(255),
    request_type    request_type NOT NULL DEFAULT 'HTTP',
    url             VARCHAR(255) NOT NULL,
    method          VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255),
    CONSTRAINT fk_action_resource_id FOREIGN KEY(resource_id) REFERENCES resources(resource_id) ON DELETE CASCADE
);

CREATE INDEX idx_actions_created_user_id ON actions(created_user_id);
CREATE INDEX idx_actions_updated_user_id ON actions(updated_user_id);
CREATE INDEX idx_actions_deleted_user_id ON actions(deleted_user_id);
CREATE INDEX idx_actions_created_at ON actions(created_at);
CREATE INDEX idx_actions_updated_at ON actions(updated_at);
CREATE INDEX idx_actions_deleted_at ON actions(deleted_at);
CREATE INDEX idx_actions_resource_id ON actions(resource_id);
