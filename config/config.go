package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port         string
	AppName      string
	StripeApiKey string
	StripeMethod string
}

func LoadConfiguration() (*Configuration, error) {
	_ = godotenv.Load()

	cfg := &Configuration{
		Port:         os.Getenv("PORT"),
		AppName:      os.Getenv("APP_NAME"),
		StripeApiKey: os.Getenv("STRIPE_API_KEY"),
		StripeMethod: os.Getenv("STRIPE_METHOD"),
	}

	if cfg.StripeApiKey == "" {
		return nil, errors.New("STRIPE_API_KEY is required")
	}

	return cfg, nil
}
