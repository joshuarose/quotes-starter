package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"gqlgen/graph/generated"
	"gqlgen/graph/model"
	"io"
	"net/http"
)

// RandomQuote is the resolver for the randomQuote field.
func (r *queryResolver) RandomQuote(ctx context.Context) (*model.Quote, error) {
	request, err := http.NewRequest("GET", "http://34.160.62.214:80/quotes", nil)
	request.Header.Set("x-api-key", "COCKTAILSAUCE")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, _ := client.Do(request)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var quote model.Quote
	json.Unmarshal(data, &quote)
	return &quote, err
}

// QuoteID is the resolver for the quoteID field.
func (r *queryResolver) QuoteID(ctx context.Context, id string) (*model.Quote, error) {
	URL := "http://34.160.62.214:80/quotes/" + id
	request, err := http.NewRequest("GET", URL, nil)
	request.Header.Set("x-api-key", "COCKTAILSAUCE")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, _ := client.Do(request)
	data, err := io.ReadAll(response.Body)
	var quote model.Quote
	json.Unmarshal(data, &quote)
	return &quote, err
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
