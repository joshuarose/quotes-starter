package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kalesecarpenter/quotes-starter/gqlgen/graph/generated"
	"github.com/kalesecarpenter/quotes-starter/gqlgen/graph/model"
)

// PostQuote is the resolver for the postQuote field.
func (r *mutationResolver) PostQuote(ctx context.Context, input model.NewQuote) (*model.Quote, error) {
	// Struct for new quote
	quote := &model.Quote{
		Quote:  input.Quote,
		Author: input.Author,
	}
	newQuote, _ := json.Marshal(&quote)
	postBody := bytes.NewBuffer(newQuote)

	// Post it to the datatbase SUCCESSFULY
	request, _ := http.NewRequest("POST", "http://34.160.48.181/quotes", postBody)
	request.Header.Set("x-api-key", "COCKTAILSAUCE")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	// This gives the ID from the created quote
	responseData, _ := io.ReadAll(response.Body)
	json.Unmarshal(responseData, &quote)
	// Send a Response Status
	return quote, nil
}

// DeleteQuote is the resolver for the deleteQuote field.
func (r *mutationResolver) DeleteQuote(ctx context.Context, id string) (*model.Quote, error) {

	//Return full quote object so user can see what is being deleted
	// Pull quote By ID from the database
	requestURL := "http://34.160.48.181/quotes/" + id
	req, _ := http.NewRequest("GET", requestURL, nil)
	// Configure the client-server connection
	client := &http.Client{}
	// throw error if not found
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	quoteByID := model.Quote{}
	quoteByIdResponse, IdResponseError := io.ReadAll(response.Body)
	if IdResponseError != nil {
		return nil, IdResponseError
	}

	// response needs to be unmarshaled into a quoted struct
	json.Unmarshal(quoteByIdResponse, &quoteByID)
	fmt.Println(quoteByID)
	// Make a request to the API for DELETING by ID
	request, _ := http.NewRequest("DELETE", "http://34.160.48.181/quotes/"+id, nil)
	// Check Header
	request.Header.Set("x-api-key", "COCKTAILSAUCE")

	deleteResponse, requestError := client.Do(request)
	if requestError != nil {
		return nil, requestError
	}

	// This gives the ID from the created quote
	_, deletedResponseError := io.ReadAll(deleteResponse.Body)
	fmt.Println(deletedResponseError)
	if deletedResponseError != nil {
		return nil, deletedResponseError
	}
	return &quoteByID, nil

}

// GetRandomQuote is the resolver for the getRandomQuote field.
func (r *queryResolver) GetRandomQuote(ctx context.Context) (*model.Quote, error) {
	// query the API endpoint
	request, err := http.NewRequest("GET", "http://34.160.48.181/quotes", nil)
	// Check Header
	request.Header.Set("x-api-key", "COCKTAILSAUCE")

	if err != nil {
		return nil, err
	}
	// Configure the client-server connection
	client := &http.Client{}
	response, _ := client.Do(request)

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var newGoQuote model.Quote
	json.Unmarshal(responseData, &newGoQuote)

	return &newGoQuote, nil
}

// GetQuoteByID is the resolver for the getQuoteByID field.
func (r *queryResolver) GetQuoteByID(ctx context.Context, id string) (*model.Quote, error) {
	// Add the ID to the end of the URL
	requestURL := "http://34.160.48.181/quotes/" + id
	request, err := http.NewRequest("GET", requestURL, nil)
	request.Header.Set("x-api-key", "COCKTAILSAUCE")

	if err != nil {
		return nil, err
	}
	// Configure the client-server connection
	client := &http.Client{}
	response, _ := client.Do(request)

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Taking the data from struct and putting it into json Quote
	var singleQuote model.Quote
	json.Unmarshal(responseData, &singleQuote)

	return &singleQuote, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
