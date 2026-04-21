CREATE TYPE "visibility_status" AS ENUM (
  'public',
  'private'
);

CREATE TYPE "q_type" AS ENUM (
  'multiple_choice',
  'true_false',
  'typed_answer'
);

CREATE TABLE "profiles" (
  "id" uuid PRIMARY KEY,
  "username" text UNIQUE NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "deleted_at" timestamp
);

CREATE TABLE "quizzes" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "author_id" uuid NOT NULL,
  "title" text NOT NULL,
  "description" text,
  "visibility" visibility_status DEFAULT 'private',
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "deleted_at" timestamp
);

CREATE TABLE "questions" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "quiz_id" uuid NOT NULL,
  "question_type" q_type NOT NULL,
  "content" text NOT NULL,
  "time_limit_seconds" int DEFAULT 15,
  "image_url" text,
  "sort_order" int NOT NULL,
  "options" jsonb NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "deleted_at" timestamp
);

CREATE TABLE "favorites" (
  "user_id" uuid,
  "quiz_id" uuid,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("user_id", "quiz_id")
);

COMMENT ON COLUMN "profiles"."id" IS 'References Supabase auth.users(id)';

COMMENT ON COLUMN "profiles"."deleted_at" IS 'GORM Soft Delete';

COMMENT ON COLUMN "quizzes"."deleted_at" IS 'GORM Soft Delete';

COMMENT ON COLUMN "questions"."sort_order" IS 'Keeps questions in order';

COMMENT ON COLUMN "questions"."options" IS 'Stores choices and correct answers';

COMMENT ON COLUMN "questions"."deleted_at" IS 'GORM Soft Delete';

ALTER TABLE "quizzes" ADD FOREIGN KEY ("author_id") REFERENCES "profiles" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "questions" ADD FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "favorites" ADD FOREIGN KEY ("user_id") REFERENCES "profiles" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "favorites" ADD FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
