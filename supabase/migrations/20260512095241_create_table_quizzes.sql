CREATE TYPE visibility_status AS ENUM ('private', 'public');

-- Create "quizzes" table
CREATE TABLE "quizzes" (
	"id" uuid NOT NULL,
	"author_id" uuid NOT NULL,
	"title" text NOT NULL,
	"description" text NULL,
	"visibility_status" visibility_status NOT NULL DEFAULT 'private',
	"created_at" timestamptz NULL DEFAULT now(),
	"updated_at" timestamptz NULL DEFAULT now(),
	"deleted_at" timestamptz NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "quizzes_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES "profiles" ("id")
);

-- Create index "idx_quizzes_deleted_at" to table: "quizzes"
CREATE INDEX "idx_quizzes_deleted_at" ON "quizzes" ("deleted_at");

-- Create index "idx_quizzes_author_id" to table: "quizzes"
CREATE INDEX "idx_quizzes_author_id" ON "quizzes" ("author_id");
