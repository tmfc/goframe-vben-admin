CREATE TABLE "sys_role_permission" (
    "role_id" bigint NOT NULL,
    "permission_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz,
    "scope" varchar(255),
    PRIMARY KEY ("role_id", "permission_id")
);
