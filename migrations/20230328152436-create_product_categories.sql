-- +migrate Up
CREATE TABLE IF NOT EXISTS product_categories (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NULL,
    lft INTEGER NOT NULL,
    rgt INTEGER NOT NULL,
    description VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(tenant_id, name),
    UNIQUE(tenant_id, lft, rgt),
    FOREIGN KEY(tenant_id) REFERENCES tenants(tenant_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE POLICY tenant_user_isolation_policy on public.product_categories
    USING (current_user = tenant_id::text);

ALTER TABLE public.product_categories ENABLE ROW LEVEL SECURITY;

-- +migrate Down
DROP TABLE IF EXISTS product_categories;