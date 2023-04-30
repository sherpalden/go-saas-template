-- +migrate Up
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ NULL
);

CREATE POLICY tenant_isolation_policy on public.tenants
    USING (current_user = tenant_id::text);

ALTER TABLE public.tenants ENABLE ROW LEVEL SECURITY;

-- +migrate Down
DROP TABLE IF EXISTS tenants;