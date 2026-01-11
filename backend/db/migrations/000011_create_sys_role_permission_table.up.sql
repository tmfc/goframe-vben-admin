CREATE TABLE "sys_role_permission" (
    "id" bigserial PRIMARY KEY,
    "role_id" bigint NOT NULL,
    "permission_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz,
    "scope" varchar(255),
    CONSTRAINT "uq_sys_role_permission_role_permission" UNIQUE ("role_id", "permission_id")
);
