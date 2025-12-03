-- Enable uuid-ossp extension needed for uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create "permissions" table
CREATE TABLE "permissions" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_permissions_deleted_at" to table: "permissions"
CREATE INDEX "idx_permissions_deleted_at" ON "permissions" ("deleted_at");
-- Create index "idx_permissions_name" to table: "permissions"
CREATE UNIQUE INDEX "idx_permissions_name" ON "permissions" ("name");
-- Create "sessions" table
CREATE TABLE "sessions" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "jwt_id" uuid NOT NULL,
  "refresh_token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked" boolean NULL DEFAULT false,
  "ip_address" text NULL,
  "user_agent" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_sessions_deleted_at" to table: "sessions"
CREATE INDEX "idx_sessions_deleted_at" ON "sessions" ("deleted_at");
-- Create index "idx_sessions_jwt_id" to table: "sessions"
CREATE UNIQUE INDEX "idx_sessions_jwt_id" ON "sessions" ("jwt_id");
-- Create index "idx_sessions_user_id" to table: "sessions"
CREATE INDEX "idx_sessions_user_id" ON "sessions" ("user_id");
-- Create "roles" table
CREATE TABLE "roles" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "roles" ("deleted_at");
-- Create index "idx_roles_name" to table: "roles"
CREATE UNIQUE INDEX "idx_roles_name" ON "roles" ("name");
-- Create "role_permissions" table
CREATE TABLE "role_permissions" (
  "role_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  CONSTRAINT "fk_permissions_roles" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_roles_permissions" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_role_permissions_permission_id" to table: "role_permissions"
CREATE INDEX "idx_role_permissions_permission_id" ON "role_permissions" ("permission_id");
-- Create index "idx_role_permissions_role_id" to table: "role_permissions"
CREATE INDEX "idx_role_permissions_role_id" ON "role_permissions" ("role_id");
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  "username" text NOT NULL,
  "email" text NOT NULL,
  "password_hash" text NOT NULL,
  "is_active" boolean NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- Create "user_profiles" table
CREATE TABLE "user_profiles" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  "phone_number" text NULL,
  "avatar_url" text NULL,
  "bio" text NULL,
  "date_of_birth" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_profile" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_user_profiles_deleted_at" to table: "user_profiles"
CREATE INDEX "idx_user_profiles_deleted_at" ON "user_profiles" ("deleted_at");
-- Create index "idx_user_profiles_user_id" to table: "user_profiles"
CREATE UNIQUE INDEX "idx_user_profiles_user_id" ON "user_profiles" ("user_id");
-- Create "user_roles" table
CREATE TABLE "user_roles" (
  "user_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  CONSTRAINT "fk_roles_users" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_users_roles" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_user_roles_role_id" to table: "user_roles"
CREATE INDEX "idx_user_roles_role_id" ON "user_roles" ("role_id");
-- Create index "idx_user_roles_user_id" to table: "user_roles"
CREATE INDEX "idx_user_roles_user_id" ON "user_roles" ("user_id");
