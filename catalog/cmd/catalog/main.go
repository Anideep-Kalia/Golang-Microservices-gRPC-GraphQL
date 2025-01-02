package main

import (
	"log"
	"time"

	"github.com/Anideep-Kalia/go-graphql-microservice/catalog"	"github.com/kelseyhightower/envconfig"
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

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("Hello ji kaise ho, aapka db nahi chal raha hai: ", err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := catalog.NewService(r)					// connecting service with db repository
	log.Fatal(catalog.ListenGRPC(s, 8080))
}