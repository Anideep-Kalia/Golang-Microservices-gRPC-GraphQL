// client will recieve the request from the server and send the request to the server and take mutations and queries commands
package account

import (
	"context"

	"github.com/github.com/Anideep-Kalia/go-graphql-grpc-micro/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn								// reference to the gRPC connection
	service pb.AccountServiceClient							// reference to the gRPC client
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())		// Insecure is used stating that the connection is not secure i.e. localhost 8080
	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn)					// new instance of the gRPC client
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(
		ctx,
		&pb.PostAccountRequest{Name: name},			// triggering the PostAccount function in the server.go file
	)
	if err != nil {
		return nil, err
	}
	return &Account{								// returing response in the form of Account struct
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(
		ctx,
		&pb.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{
			Skip: skip,
			Take: take,
		},
	)
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}
	return accounts, nil
}