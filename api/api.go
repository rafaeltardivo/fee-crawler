package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/rafaeltardivo/fee-crawler/crawler"
	"github.com/rafaeltardivo/fee-crawler/rates"
)

// Fee type definition
type feeStruct struct {
	RatesDate   string `json:"rates_date"`
	Description string `json:"description,omitempty"`
	BRL         string `json:"BRL"`
	USD         string `json:"USD"`
	EUR         string `json:"EUR"`
}

// Fee GraphQL type definition
var feeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Fee",
	Fields: graphql.Fields{
		"rates_date": &graphql.Field{
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

// Fee GraphQL root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"transfer": &graphql.Field{
			Type: feeType,
			Args: graphql.FieldConfigArgument{
				"domain": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"plan": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: feeResolver,
		},
	},
})

// Fee GraphQL resolver
// Parallelizes requests to  rates and crawler services and return results
func feeResolver(p graphql.ResolveParams) (interface{}, error) {
	domain := p.Args["domain"].(string)
	plan := p.Args["plan"].(string)

	logger.Info(fmt.Sprintf("processing request with domain: %s plan: %s", domain, plan))
	if !validateDomain(domain) {
		return nil, apiError(fmt.Sprintf("invalid domain %s", domain))
	}

	var wg sync.WaitGroup
	wg.Add(2)

	crawlerChannel := make(chan crawler.CrawlerResponse)
	ratesChannel := make(chan rates.RatesResponse)

	go crawler.CrawlFeeData(domain, plan, crawlerChannel, &wg)
	crawlerResponse := <-crawlerChannel

	go rates.GetRatesData(ratesChannel, &wg)
	ratesResponse := <-ratesChannel
	wg.Wait()

	return toAPIResponse(ratesResponse, crawlerResponse)
}

// Serves GraphQL API
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
