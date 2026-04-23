-- update own profile
DROP POLICY IF EXISTS "users can update own profile" ON public.profiles;
CREATE POLICY "users can update own profile"
  ON public.profiles FOR UPDATE
  USING ((select auth.uid()) = id);

-- drop delete own profile  
DROP POLICY IF EXISTS "users can delete own profile" ON public.profiles;