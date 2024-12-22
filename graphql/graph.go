// this file is like graphql gateway which connects to different micro-services and provides the resolvers for the graphql schema
package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/Anideep-Kalia/go-graphql-grpc-micro/account"
	"github.com/Anideep-Kalia/go-graphql-grpc-micro/catalog"
	"github.com/Anideep-Kalia/go-graphql-grpc-micro/order"
)

// struct uses to connect to the different micro-services
type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogURL, orderURL string) (*Server, error) {
	// Connect to account service
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	// Connect to product service
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	// Connect to order service
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

// generates a schema with mappings to the resolvers
func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}