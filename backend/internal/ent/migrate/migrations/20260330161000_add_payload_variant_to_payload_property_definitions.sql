-- Modify "payload_property_definitions" table
ALTER TABLE "payload_property_definitions" ADD COLUMN "payload_variant" character varying NOT NULL DEFAULT '';

-- Drop old indexes
DROP INDEX IF EXISTS "payloadpropertydefinition_payload_type_key";
DROP INDEX IF EXISTS "payloadpropertydefinition_payload_type_order_index";
DROP INDEX IF EXISTS "payloadpropertydefinition_payload_type_is_nested";

-- Create new indexes with payload_variant support
CREATE UNIQUE INDEX "payloadpropertydefinition_payload_type_payload_variant_key" ON "payload_property_definitions" ("payload_type", "payload_variant", "key");
CREATE INDEX "payloadpropertydefinition_payload_type_payload_variant_order_index" ON "payload_property_definitions" ("payload_type", "payload_variant", "order_index");
CREATE INDEX "payloadpropertydefinition_payload_type_payload_variant_is_nested" ON "payload_property_definitions" ("payload_type", "payload_variant", "is_nested");
CREATE INDEX "payloadpropertydefinition_payload_type" ON "payload_property_definitions" ("payload_type");
