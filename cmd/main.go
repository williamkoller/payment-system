package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/williamkoller/payment-system/config"
	healthRouter "github.com/williamkoller/payment-system/internal/healthz/router"
	"github.com/williamkoller/payment-system/internal/middleware"
	paymentRouter "github.com/williamkoller/payment-system/internal/payment/router"
	"github.com/williamkoller/payment-system/pkg/logger"
)

func main() {
	r := gin.Default()

	logger.InitLogger()
	defer logger.Sync()

	configuration, err := config.LoadConfiguration()

	if err != nil {
		log.Fatal(err)
	}

	middleware.Middlewares(r)
	healthRouter.SetupRouter(r)
	paymentRouter.SetupRouter(r)

	srv := &http.Server{
		Addr:              ":" + configuration.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("Starting server", "AppName", configuration.AppName)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalw("Error starting server", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	logger.Info("Server shutting down")
}
