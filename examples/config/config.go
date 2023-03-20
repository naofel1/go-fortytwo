package config

import (
	"context"
	"log"

	"github.com/naofel1/go-fortytwo"
)

// API42 struct contains all the information needed to initialize the 42 API client
type API42 struct {
	ClientID     string   `json:"client_id" yaml:"client_id"`
	ClientSecret string   `json:"client_secret" yaml:"client_secret"`
	RedirectURL  string   `json:"redirect_url" yaml:"redirect_url"`
	Scopes       []string `json:"scopes" yaml:"scopes"`
}

// Init42API will return a 42 client initialized
func Init42API(ctx context.Context, cfg *API42) *fortytwo.Client {
	client, err := fortytwo.NewClient(ctx, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, cfg.Scopes)
	if err != nil {
		log.Fatal("client initialization failed", err)
	}

	log.Printf("42API client initialized")

	return client
}
