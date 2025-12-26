CREATE TABLE "sys_permission" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "description" text,
    "parent_id" bigint,
    "status" int,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz,
    "creator_id" bigint,
    "modifier_id" bigint,
    "dept_id" bigint
);
