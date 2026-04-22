-- Create "profiles" table
CREATE TABLE "profiles" (
  "id" uuid NOT NULL,
  "username" text NOT NULL,
  "full_name" text NULL,
  "avatar_url" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_profiles_username" UNIQUE ("username")
);
-- Create index "idx_profiles_deleted_at" to table: "profiles"
CREATE INDEX "idx_profiles_deleted_at" ON "profiles" ("deleted_at");
