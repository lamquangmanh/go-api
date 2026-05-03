-- Create modules table
CREATE TABLE modules (
    module_id       UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    description     VARCHAR(255),
    url             VARCHAR(255),
    icon            VARCHAR(255),
    product_id      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255)
);

CREATE INDEX idx_modules_created_user_id ON modules(created_user_id);
CREATE INDEX idx_modules_updated_user_id ON modules(updated_user_id);
CREATE INDEX idx_modules_deleted_user_id ON modules(deleted_user_id);
CREATE INDEX idx_modules_created_at ON modules(created_at);
CREATE INDEX idx_modules_updated_at ON modules(updated_at);
CREATE INDEX idx_modules_deleted_at ON modules(deleted_at);
CREATE INDEX idx_modules_product_id ON modules(product_id);
