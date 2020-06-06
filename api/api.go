package api

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
)

type fee struct {
	CalculatedAt string `json:"calculated_at"`
	Description  string `json:"description,omitempty"`
	BRL          string `json:"BRL"`
	USD          string `json:"USD"`
	EUR          string `json:"EUR"`
}

var feeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Fee",
	Fields: graphql.Fields{
		"calculated_at": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"BRL": &graphql.Field{
			Type: graphql.String,
		},
		"USD": &graphql.Field{
			Type: graphql.String,
		},
		"EUR": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"domain": &graphql.Field{
			Type: feeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return fee{
					CalculatedAt: "2020-03-03",
					Description:  "Bla Bla Bla",
					BRL:          "1.50",
					USD:          "2.00",
					EUR:          "3.00",
				}, nil
			},
		},
	},
})

func Serve(w http.ResponseWriter, r *http.Request) {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: r.URL.Query().Get("query"),
	})
	json.NewEncoder(w).Encode(result)
}
