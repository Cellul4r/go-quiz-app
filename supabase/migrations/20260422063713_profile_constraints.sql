-- Link profiles.id to auth.users.id
ALTER TABLE "profiles"
  ADD CONSTRAINT "profiles_id_fkey"
  FOREIGN KEY ("id")
  REFERENCES auth.users(id)
  ON DELETE CASCADE;

-- Enable RLS
ALTER TABLE "profiles" ENABLE ROW LEVEL SECURITY;

-- RLS policies
CREATE POLICY "public profiles are viewable by everyone"
  ON "profiles" FOR SELECT
  USING (true);

CREATE POLICY "users can update own profile"
  ON "profiles" FOR UPDATE
  USING (auth.uid() = id);

CREATE POLICY "users can delete own profile"
  ON "profiles" FOR DELETE
  USING (auth.uid() = id);