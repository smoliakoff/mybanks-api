.PHONY: generate enrich openapi genroutes

generate:
	@echo "🔧 Generating Ent client + OpenAPI JSON spec..."
	go generate ./...
	go generate ./ent
	go run github.com/99designs/gqlgen generate

enrich: generate
	@echo "➕ Injecting metadata into openapi.json..."
	yq -i ' \
  .info = { \
    "title": "Bank Directory API", \
    "description": "API for bank directory", \
    "version": "1.0.0" \
  } | \
  .servers = [{ \
    "url": "http://127.0.0.1:8080/api/v1", \
    "description": "Local server (default)" \
  }]' ent/openapi.json

	@echo "✅ openapi.json enriched"

openapi: enrich
	@echo "🔁 Converting enriched JSON to YAML..."
	yq -P -o=yaml ent/openapi.json > api/openapi.yaml
	@echo "✅ docs/openapi.generated.yaml created from enriched JSON"

