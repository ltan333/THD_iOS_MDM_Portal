-- Modify "payload_property_definitions" table
ALTER TABLE "payload_property_definitions" ALTER COLUMN "description" TYPE text, ALTER COLUMN "order_index" TYPE bigint, ADD COLUMN "title" character varying NULL, ADD COLUMN "presence" character varying NOT NULL DEFAULT 'optional', ADD COLUMN "supported_os" jsonb NULL, ADD COLUMN "conditions" jsonb NULL, ADD COLUMN "yaml_source_file" character varying NULL;
-- Create index "payloadpropertydefinition_payload_type_is_nested" to table: "payload_property_definitions"
CREATE INDEX "payloadpropertydefinition_payload_type_is_nested" ON "payload_property_definitions" ("payload_type", "is_nested");
-- Create index "payloadpropertydefinition_payload_type_order_index" to table: "payload_property_definitions"
CREATE INDEX "payloadpropertydefinition_payload_type_order_index" ON "payload_property_definitions" ("payload_type", "order_index");
