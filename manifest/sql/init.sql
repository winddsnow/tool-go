-- User table
CREATE TABLE IF NOT EXISTS "user" (
    "id" BIGSERIAL PRIMARY KEY,
    "username" VARCHAR(64) NOT NULL UNIQUE,
    "password" VARCHAR(128) NOT NULL,
    "salt" VARCHAR(32) NOT NULL,
    "nickname" VARCHAR(64) DEFAULT '',
    "email" VARCHAR(128) DEFAULT '',
    "phone" VARCHAR(32) DEFAULT '',
    "status" SMALLINT NOT NULL DEFAULT 1,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_user_username" ON "user" ("username");
CREATE INDEX "idx_user_status" ON "user" ("status");
CREATE INDEX "idx_user_deleted_at" ON "user" ("deleted_at");

COMMENT ON TABLE "user" IS '用户表';
COMMENT ON COLUMN "user"."id" IS '用户ID';
COMMENT ON COLUMN "user"."username" IS '用户名';
COMMENT ON COLUMN "user"."password" IS '密码哈希 (MD5 + Salt)';
COMMENT ON COLUMN "user"."salt" IS '密码盐值';
COMMENT ON COLUMN "user"."nickname" IS '昵称';
COMMENT ON COLUMN "user"."email" IS '邮箱';
COMMENT ON COLUMN "user"."phone" IS '手机号';
COMMENT ON COLUMN "user"."status" IS '状态: 1=启用, 0=禁用';

-- Role table
CREATE TABLE IF NOT EXISTS "role" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(64) NOT NULL,
    "code" VARCHAR(64) NOT NULL UNIQUE,
    "sort" INT NOT NULL DEFAULT 0,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "desc" VARCHAR(256) DEFAULT '',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_role_code" ON "role" ("code");
CREATE INDEX "idx_role_status" ON "role" ("status");
CREATE INDEX "idx_role_deleted_at" ON "role" ("deleted_at");

COMMENT ON TABLE "role" IS '角色表';
COMMENT ON COLUMN "role"."id" IS '角色ID';
COMMENT ON COLUMN "role"."name" IS '角色名称';
COMMENT ON COLUMN "role"."code" IS '角色编码';
COMMENT ON COLUMN "role"."sort" IS '排序';
COMMENT ON COLUMN "role"."status" IS '状态: 1=启用, 0=禁用';
COMMENT ON COLUMN "role"."desc" IS '描述';

-- User-Role relation table
CREATE TABLE IF NOT EXISTS "user_role" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL,
    "role_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX "idx_user_role_user_role" ON "user_role" ("user_id", "role_id");

COMMENT ON TABLE "user_role" IS '用户角色关联表';

-- 默认管理员账号 (密码: walter)
INSERT INTO "user" ("username", "password", "salt", "nickname", "status") VALUES
('walter', 'd1c2ae9977ade915fb4e507c40eb40b6', '63da31ad9cb14c63c887433e21b06b21', '本地开发', 1);

-- 默认角色
INSERT INTO "role" ("name", "code", "sort", "status", "desc") VALUES
('超级管理员', 'super_admin', 1, 1, '拥有所有权限'),
('管理员', 'admin', 2, 1, '管理员'),
('普通用户', 'user', 3, 1, '普通用户');

-- 分配管理员角色
INSERT INTO "user_role" ("user_id", "role_id") VALUES (1, 1);

-- Page view tracking table
CREATE TABLE IF NOT EXISTS "page_view" (
    "id" BIGSERIAL PRIMARY KEY,
    "page_path" VARCHAR(255) NOT NULL DEFAULT '',
    "user_id" BIGINT DEFAULT NULL,
    "username" VARCHAR(64) DEFAULT '',
    "ip_address" VARCHAR(64) DEFAULT '',
    "user_agent" VARCHAR(512) DEFAULT '',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_page_view_user_id" ON "page_view" ("user_id");
CREATE INDEX "idx_page_view_created_at" ON "page_view" ("created_at");

COMMENT ON TABLE "page_view" IS '页面访问埋点表';
COMMENT ON COLUMN "page_view"."page_path" IS '访问页面路径';
COMMENT ON COLUMN "page_view"."user_id" IS '用户ID (未登录为空)';
COMMENT ON COLUMN "page_view"."username" IS '用户名';
COMMENT ON COLUMN "page_view"."ip_address" IS 'IP地址';
COMMENT ON COLUMN "page_view"."user_agent" IS 'User-Agent';
