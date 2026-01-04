CREATE TABLE sys_dept (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES sys_tenant(id),
    parent_id UUID,
    name VARCHAR(255) NOT NULL,
    "order" INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    creator_id BIGINT,
    modifier_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_sys_dept_tenant ON sys_dept (tenant_id);
CREATE INDEX idx_sys_dept_parent ON sys_dept (parent_id);
CREATE INDEX idx_sys_dept_status ON sys_dept (status);