-- Create "favorites" table
CREATE TABLE "favorites" (
	"profile_id" uuid NOT NULL,
	"quiz_id" uuid NOT NULL,
	"created_at" timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY ("profile_id", "quiz_id"),
	CONSTRAINT "favorites_profile_id_fkey" FOREIGN KEY ("profile_id") REFERENCES "profiles" ("id") ON DELETE CASCADE,
	CONSTRAINT "favorites_quiz_id_fkey" FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id") ON DELETE CASCADE
);

-- Create index "idx_favorites_profile_id_created_at" to table: "favorites"
CREATE INDEX "idx_favorites_profile_id_created_at" ON "favorites" ("profile_id", "created_at" DESC);

-- Create index "idx_favorites_quiz_id" to table: "favorites"
CREATE INDEX "idx_favorites_quiz_id" ON "favorites" ("quiz_id");
