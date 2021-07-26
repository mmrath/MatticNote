package graphql

import (
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_mutation"
	"github.com/graphql-go/graphql"
)

var mutationRoot = graphql.ObjectConfig{
	Name:        "MNMutation",
	Description: "MatticNote Mutation",
	Fields: graphql.Fields{
		"createNote": mn_mutation.CreateNote,
		"deleteNote": mn_mutation.DeleteNote,
		"followUser": mn_mutation.FollowUser,
	},
}
