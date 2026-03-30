-- Add items_type, items_reference, is_nested, order_index columns to payload_property_definitions
-- v1.2: support array-of-dictionary metadata and hierarchy helpers

ALTER TABLE "payload_property_definitions"
  ADD COLUMN IF NOT EXISTS "items_type" character varying NULL,
  ADD COLUMN IF NOT EXISTS "items_reference" character varying NULL,
  ADD COLUMN IF NOT EXISTS "is_nested" boolean NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS "order_index" integer NOT NULL DEFAULT 0;
