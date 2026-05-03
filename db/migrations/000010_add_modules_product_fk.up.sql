-- Add foreign key constraint modules.product_id -> products.product_id
ALTER TABLE modules
ADD CONSTRAINT fk_modules_product_id
    FOREIGN KEY (product_id)
    REFERENCES products(product_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;
