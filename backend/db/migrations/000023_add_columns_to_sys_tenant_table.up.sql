ALTER TABLE sys_tenant ADD COLUMN code VARCHAR(50);
ALTER TABLE sys_tenant ADD COLUMN contact_name VARCHAR(100);
ALTER TABLE sys_tenant ADD COLUMN contact_mobile VARCHAR(20);
ALTER TABLE sys_tenant ADD COLUMN start_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE sys_tenant ADD COLUMN expire_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE sys_tenant ADD COLUMN package_version VARCHAR(50);
ALTER TABLE sys_tenant ADD COLUMN domain VARCHAR(255);

CREATE UNIQUE INDEX idx_sys_tenant_code ON sys_tenant(code);
