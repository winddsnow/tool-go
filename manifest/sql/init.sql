-- User table
CREATE TABLE IF NOT EXISTS "user" (
    "id" BIGSERIAL PRIMARY KEY,
    "username" VARCHAR(64) NOT NULL UNIQUE,
    "password" VARCHAR(128) NOT NULL,
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

COMMENT ON TABLE "user" IS 'User table';
COMMENT ON COLUMN "user"."id" IS 'User ID';
COMMENT ON COLUMN "user"."username" IS 'Username';
COMMENT ON COLUMN "user"."password" IS 'Password';
COMMENT ON COLUMN "user"."nickname" IS 'Nickname';
COMMENT ON COLUMN "user"."email" IS 'Email';
COMMENT ON COLUMN "user"."phone" IS 'Phone';
COMMENT ON COLUMN "user"."status" IS 'Status: 1=active, 0=disabled';

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

COMMENT ON TABLE "role" IS 'Role table';
COMMENT ON COLUMN "role"."id" IS 'Role ID';
COMMENT ON COLUMN "role"."name" IS 'Role name';
COMMENT ON COLUMN "role"."code" IS 'Role code';
COMMENT ON COLUMN "role"."sort" IS 'Sort order';
COMMENT ON COLUMN "role"."status" IS 'Status: 1=active, 0=disabled';
COMMENT ON COLUMN "role"."desc" IS 'Description';

-- User-Role relation table
CREATE TABLE IF NOT EXISTS "user_role" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL,
    "role_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX "idx_user_role_user_role" ON "user_role" ("user_id", "role_id");

COMMENT ON TABLE "user_role" IS 'User-Role relation table';

-- Insert default admin user (password: admin123)
INSERT INTO "user" ("username", "password", "nickname", "status") VALUES
('admin', 'admin123', 'Administrator', 1);

-- Insert default roles
INSERT INTO "role" ("name", "code", "sort", "status", "desc") VALUES
('Super Admin', 'super_admin', 1, 1, 'Super administrator with all permissions'),
('Admin', 'admin', 2, 1, 'Administrator'),
('User', 'user', 3, 1, 'Regular user');

-- Assign admin role to admin user
INSERT INTO "user_role" ("user_id", "role_id") VALUES (1, 1);
