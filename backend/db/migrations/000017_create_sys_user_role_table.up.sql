CREATE TABLE "sys_user_role" (
    "id" BIGSERIAL PRIMARY KEY,
    "tenant_id" BIGINT NOT NULL REFERENCES sys_tenant(id),
    "user_id" BIGINT NOT NULL REFERENCES sys_user(id) ON DELETE CASCADE,
    "role_id" BIGINT NOT NULL REFERENCES sys_role(id) ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "created_by" BIGINT,
    CONSTRAINT "uq_sys_user_role" UNIQUE ("tenant_id", "user_id", "role_id")
);

CREATE INDEX "idx_sys_user_role_user_id" ON "sys_user_role" ("user_id");
CREATE INDEX "idx_sys_user_role_role_id" ON "sys_user_role" ("role_id");
CREATE INDEX "idx_sys_user_role_tenant_id" ON "sys_user_role" ("tenant_id");
