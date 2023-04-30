-- +migrate Up
CREATE TABLE IF NOT EXISTS tenants_users (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(255) NOT NULL,
    account_status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ NULL,
    UNIQUE(tenant_id, user_id),
    FOREIGN KEY(tenant_id) REFERENCES tenants(tenant_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE POLICY tenant_user_isolation_policy on public.tenants_users
    USING (current_user = tenant_id::text);

ALTER TABLE public.tenants_users ENABLE ROW LEVEL SECURITY;

-- +migrate Down
DROP TABLE IF EXISTS tenants_users;