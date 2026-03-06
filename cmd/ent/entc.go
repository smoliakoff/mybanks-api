package main

import (
	"errors"
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaPath("../graph/ent.graphqls"),
		entgql.WithSchemaGenerator(), // enable GraphQL schema generate
		entgql.WithWhereInputs(true), // with filters
	)
	if !errors.Is(err, nil) {
		log.Fatalf("Error: failed creating entgql extension: %v", err)
	}
	if err := entc.Generate("./schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureGlobalID,
		},
	}, entc.Extensions(ex)); !errors.Is(err, nil) {
		log.Fatalf("Error: failed running ent codegen: %v", err)
	}
}
