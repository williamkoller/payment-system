package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

func LoadConfiguration() (*ResponseConfiguration, error) {
	app, err := loadAppConfiguration()
	if err != nil {
		return nil, fmt.Errorf("Error loading app configuration: %w", err)
	}

	stripe, err := loadStripeConfiguration()
	if err != nil {
		return nil, fmt.Errorf("Error loading Stripe configuration: %w", err)
	}

	db, err := LoadDatabaseConfiguration()
	if err != nil {
		return nil, fmt.Errorf("Error loading database configuration: %w", err)
	}

	return &ResponseConfiguration{
		App:      *app,
		Stripe:   *stripe,
		Database: *db,
	}, nil
}

func loadAppConfiguration() (*AppConfiguration, error) {
	app := &AppConfiguration{
		Port:    os.Getenv("PORT"),
		AppName: os.Getenv("APP_NAME"),
	}
	if app.Port == "" {
		return nil, errors.New("PORT is required")
	}
	if app.AppName == "" {
		return nil, errors.New("APP_NAME is required")
	}
	return app, nil
}
