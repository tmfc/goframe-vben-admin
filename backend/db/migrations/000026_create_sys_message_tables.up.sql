CREATE TABLE sys_message (
    id          BIGSERIAL PRIMARY KEY,
    tenant_id   BIGINT NOT NULL REFERENCES sys_tenant(id),
    title       VARCHAR(255) NOT NULL,
    content     TEXT,
    type        SMALLINT NOT NULL DEFAULT 1,
    receiver_id BIGINT,
    path        VARCHAR(255),
    params      JSONB,
    status      SMALLINT DEFAULT 1,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE,
    creator_id  BIGINT
);

CREATE INDEX idx_sys_message_tenant ON sys_message (tenant_id);
CREATE INDEX idx_sys_message_receiver ON sys_message (receiver_id);
CREATE INDEX idx_sys_message_type ON sys_message (type);

CREATE TABLE sys_message_read (
    id          BIGSERIAL PRIMARY KEY,
    message_id  BIGINT NOT NULL REFERENCES sys_message(id) ON DELETE CASCADE,
    user_id     BIGINT NOT NULL,
    read_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sys_message_read_message_user ON sys_message_read (message_id, user_id);
CREATE INDEX idx_sys_message_read_user ON sys_message_read (user_id);
