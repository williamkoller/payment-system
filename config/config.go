package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfiguration struct {
	Port    string
	AppName string
}

type StripeConfiguration struct {
	StripeApiKey  string
	StripeMethod  string
	StripeWebhook string
}

type DatabaseConfiguration struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type ResponseConfiguration struct {
	App      AppConfiguration
	Stripe   StripeConfiguration
	Database DatabaseConfiguration
}

func LoadConfiguration() (*ResponseConfiguration, error) {
	_ = godotenv.Load()

	app, err := loadAppConfiguration()
	if err != nil {
		return nil, errors.New("Error loading app configuration: " + err.Error())
	}

	stripe, err := loadStripeConfiguration()

	if err != nil {
		return nil, errors.New("Error loading Stripe configuration: " + err.Error())
	}

	return &ResponseConfiguration{
		App:    *app,
		Stripe: *stripe,
	}, nil
}

func loadAppConfiguration() (*AppConfiguration, error) {
	app := &AppConfiguration{
		Port:    os.Getenv("PORT"),
		AppName: os.Getenv("APP_NAME"),
	}

	return app, nil
}

func loadStripeConfiguration() (*StripeConfiguration, error) {
	stripe := &StripeConfiguration{
		StripeApiKey:  os.Getenv("STRIPE_API_KEY"),
		StripeMethod:  os.Getenv("STRIPE_METHOD"),
		StripeWebhook: os.Getenv("STRIPE_WEBHOOK"),
	}

	if stripe.StripeApiKey == "" {
		return nil, errors.New("STRIPE_API_KEY is required")
	}

	return stripe, nil

}

func LoadDatabaseConfiguration() (*DatabaseConfiguration, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}
	db := &DatabaseConfiguration{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}

	return db, nil
}
