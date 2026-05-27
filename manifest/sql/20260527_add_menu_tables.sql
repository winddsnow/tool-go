-- Migration: Add menu + role_menu tables
-- Date: 2026-05-27

CREATE TABLE IF NOT EXISTS "menu" (
    "id" BIGSERIAL PRIMARY KEY,
    "parent_id" BIGINT NOT NULL DEFAULT 0,
    "name" VARCHAR(64) NOT NULL,
    "path" VARCHAR(128) NOT NULL DEFAULT '',
    "component" VARCHAR(256) NOT NULL DEFAULT '',
    "icon" VARCHAR(64) NOT NULL DEFAULT '',
    "sort" INT NOT NULL DEFAULT 0,
    "visible" SMALLINT NOT NULL DEFAULT 1,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "type" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS "idx_menu_parent_id" ON "menu" ("parent_id");
CREATE INDEX IF NOT EXISTS "idx_menu_status" ON "menu" ("status");
CREATE INDEX IF NOT EXISTS "idx_menu_deleted_at" ON "menu" ("deleted_at");

COMMENT ON TABLE "menu" IS '菜单表';
COMMENT ON COLUMN "menu"."id" IS '菜单ID';
COMMENT ON COLUMN "menu"."parent_id" IS '父菜单ID (0=顶级)';
COMMENT ON COLUMN "menu"."name" IS '菜单名称';
COMMENT ON COLUMN "menu"."path" IS '路由路径';
COMMENT ON COLUMN "menu"."component" IS '前端组件路径';
COMMENT ON COLUMN "menu"."icon" IS '图标名';
COMMENT ON COLUMN "menu"."sort" IS '排序';
COMMENT ON COLUMN "menu"."visible" IS '是否显示: 1=显示, 0=隐藏';
COMMENT ON COLUMN "menu"."status" IS '状态: 1=启用, 0=禁用';
COMMENT ON COLUMN "menu"."type" IS '类型: 1=目录, 2=菜单, 3=按钮';

CREATE TABLE IF NOT EXISTS "role_menu" (
    "id" BIGSERIAL PRIMARY KEY,
    "role_id" BIGINT NOT NULL,
    "menu_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_role_menu_role_menu" ON "role_menu" ("role_id", "menu_id");
COMMENT ON TABLE "role_menu" IS '角色菜单关联表';

-- Seed menus
INSERT INTO "menu" ("id", "parent_id", "name", "path", "component", "icon", "sort", "visible", "status", "type") VALUES
(1, 0, '工具箱',   '/tools',       'views/tools/index.vue',          'Tool',      1, 1, 1, 1),
(2, 0, '系统管理', '/system',      '',                              'Setting',   2, 1, 1, 1),
(3, 2, '用户管理', '/system/user', 'views/system/user/index.vue',   'User',      1, 1, 1, 2),
(4, 2, '角色管理', '/system/role', 'views/system/role/index.vue',   'Avatar',    2, 1, 1, 2),
(5, 2, '菜单管理', '/system/menu', 'views/system/menu/index.vue',   'Menu',      3, 1, 1, 2)
ON CONFLICT ("id") DO NOTHING;

-- Seed role_menu
INSERT INTO "role_menu" ("role_id", "menu_id") VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
(2, 1), (2, 2), (2, 3), (2, 4),
(3, 1)
ON CONFLICT ("role_id", "menu_id") DO NOTHING;
