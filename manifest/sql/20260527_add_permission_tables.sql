-- Migration: Add permission + role_permission tables
-- Date: 2026-05-27

CREATE TABLE IF NOT EXISTS "permission" (
    "id" BIGSERIAL PRIMARY KEY,
    "code" VARCHAR(64) NOT NULL,
    "name" VARCHAR(64) NOT NULL,
    "menu_id" BIGINT NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_permission_code" ON "permission" ("code");

COMMENT ON TABLE "permission" IS '权限表';
COMMENT ON COLUMN "permission"."id" IS '权限ID';
COMMENT ON COLUMN "permission"."code" IS '权限码 (resource:action)';
COMMENT ON COLUMN "permission"."name" IS '权限名称';
COMMENT ON COLUMN "permission"."menu_id" IS '关联菜单ID (0=无关联)';

CREATE TABLE IF NOT EXISTS "role_permission" (
    "id" BIGSERIAL PRIMARY KEY,
    "role_id" BIGINT NOT NULL,
    "permission_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_role_permission_role_perm" ON "role_permission" ("role_id", "permission_id");

COMMENT ON TABLE "role_permission" IS '角色权限关联表';

-- Seed permissions
INSERT INTO "permission" ("code", "name", "menu_id") VALUES
('user:create',       '创建用户',   3),
('user:delete',       '删除用户',   3),
('user:assign-roles', '分配角色',   3),
('role:create',       '创建角色',   4),
('role:delete',       '删除角色',   4),
('menu:create',       '创建菜单',   5),
('menu:delete',       '删除菜单',   5);

-- Seed role_permission
-- super_admin: all permissions
INSERT INTO "role_permission" ("role_id", "permission_id") VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7);

-- admin: user CRUD + role create/delete (no menu management)
INSERT INTO "role_permission" ("role_id", "permission_id") VALUES
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

-- user: no button-level permissions
