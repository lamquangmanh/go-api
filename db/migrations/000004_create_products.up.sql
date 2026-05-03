-- Create products table
CREATE TABLE products (
    product_id      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    description     VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_user_id VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_user_id VARCHAR(255),
    deleted_at      TIMESTAMPTZ,
    deleted_user_id VARCHAR(255)
);

CREATE INDEX idx_products_created_user_id ON products(created_user_id);
CREATE INDEX idx_products_updated_user_id ON products(updated_user_id);
CREATE INDEX idx_products_deleted_user_id ON products(deleted_user_id);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_updated_at ON products(updated_at);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
