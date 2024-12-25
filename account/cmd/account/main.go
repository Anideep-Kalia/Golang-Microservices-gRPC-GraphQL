package main

import (
	"log"
	"time"

	"github.com/github.com/Anideep-Kalia/go-graphql-grpc-micro/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {				// retrying until connectected with db and retry after 2 seconds
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)				// NewPostgresRepository is a function in account/repository.go
		if err != nil {
			log.Println(err)												// this will indicate that the connection is not established and to retry
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := account.NewService(r)						     					// NewService is a function in service.go		
	log.Fatal(account.ListenGRPC(s, 8080))									// ListenGRPC is a function in server.go
}