package main

import (
	"context"
	"os"

	"github.com/Udehlee/alert-Me/api"
	"github.com/Udehlee/alert-Me/db/db"
	"github.com/Udehlee/alert-Me/pkg/rabbitmq"
	"github.com/Udehlee/alert-Me/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
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

	rb := rabbitmq.NewRabbitMQ(rbConn.Conn, rbConn.Ch, dbConn)

	svc := service.NewService(rb)
	svc.StartConsumer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go svc.PeriodicCheck(ctx, "product-check")

	h := api.NewHandler(log, rb)
	h.RegisterRoutes(r)

	if err := r.Run(":8000"); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}

}
