package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Udehlee/alert-Me/api"
	db "github.com/Udehlee/alert-Me/internals/db/conn"
	"github.com/Udehlee/alert-Me/internals/rabbitmq"
	"github.com/Udehlee/alert-Me/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	log := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()

	dbConn, err := db.InitConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	rbConn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to RabbitMQ")
	}

	rb := rabbitmq.NewRabbitMQ(rbConn.Conn, rbConn.Ch)

	svc := service.NewService(dbConn, rb)
	svc.StartConsumer()
	svc.ComparePrice("product_check")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go svc.SendForRecheck(ctx, "product_check")

	h := api.NewHandler(log, *svc)
	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		log.Info().Msg("Starting HTTP server on :8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	if err := srv.Shutdown(ctx); err != nil {
		log.Err(err).Msg("server stopping forcefully")
	}

	cancel()
	log.Info().Msg("Server stopped")
}
