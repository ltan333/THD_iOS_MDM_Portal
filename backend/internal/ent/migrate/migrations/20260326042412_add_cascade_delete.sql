-- Modify "profile_assignments" table
ALTER TABLE "public"."profile_assignments" DROP CONSTRAINT "profile_assignments_profiles_assignments", ADD CONSTRAINT "profile_assignments_profiles_assignments" FOREIGN KEY ("profile_id") REFERENCES "public"."profiles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "profile_deployment_statuses" table
ALTER TABLE "public"."profile_deployment_statuses" DROP CONSTRAINT "profile_deployment_statuses_profiles_deployment_statuses", ADD CONSTRAINT "profile_deployment_statuses_profiles_deployment_statuses" FOREIGN KEY ("profile_id") REFERENCES "public"."profiles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "profile_versions" table
ALTER TABLE "public"."profile_versions" DROP CONSTRAINT "profile_versions_profiles_versions", ADD CONSTRAINT "profile_versions_profiles_versions" FOREIGN KEY ("profile_id") REFERENCES "public"."profiles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
