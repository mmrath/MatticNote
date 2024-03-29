package graphql

import "github.com/graphql-go/graphql"

var schemaConfig = graphql.SchemaConfig{
	Query:    graphql.NewObject(queryRoot),
	Mutation: graphql.NewObject(mutationRoot),
}

func setupSchema() graphql.Schema {
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}
	return schema
}
