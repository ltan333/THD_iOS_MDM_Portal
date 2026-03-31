-- Create "dep_devices" table
CREATE TABLE "dep_devices" (
  "id" character varying NOT NULL,
  "serial_number" character varying NOT NULL,
  "dep_name" character varying NOT NULL,
  "model" character varying NULL,
  "description" character varying NULL,
  "color" character varying NULL,
  "asset_tag" character varying NULL,
  "os" character varying NULL,
  "device_family" character varying NULL,
  "profile_uuid" character varying NULL,
  "profile_status" character varying NULL,
  "profile_assign_time" timestamptz NULL,
  "profile_push_time" timestamptz NULL,
  "device_assigned_by" character varying NULL,
  "device_assigned_date" timestamptz NULL,
  "op_type" character varying NULL,
  "op_date" timestamptz NULL,
  "needs_manual_reassign" boolean NOT NULL DEFAULT false,
  "reassign_error" character varying NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "depdevice_dep_name" to table: "dep_devices"
CREATE INDEX "depdevice_dep_name" ON "dep_devices" ("dep_name");
-- Create index "depdevice_needs_manual_reassign" to table: "dep_devices"
CREATE INDEX "depdevice_needs_manual_reassign" ON "dep_devices" ("needs_manual_reassign");
-- Create index "depdevice_op_type" to table: "dep_devices"
CREATE INDEX "depdevice_op_type" ON "dep_devices" ("op_type");
-- Create index "depdevice_profile_uuid" to table: "dep_devices"
CREATE INDEX "depdevice_profile_uuid" ON "dep_devices" ("profile_uuid");
-- Create index "depdevice_serial_number" to table: "dep_devices"
CREATE UNIQUE INDEX "depdevice_serial_number" ON "dep_devices" ("serial_number");
