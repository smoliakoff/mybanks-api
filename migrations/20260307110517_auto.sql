-- Modify "currency_rates" table
ALTER TABLE "public"."currency_rates" ADD COLUMN "base" character varying NOT NULL DEFAULT 'GEL';
-- Create index "currencyrate_base_currency" to table: "currency_rates"
CREATE INDEX "currencyrate_base_currency" ON "public"."currency_rates" ("base", "currency");
