package main

import (
	"context"
	"os"

	"github.com/Udehlee/alert-Me/api"
	"github.com/Udehlee/alert-Me/internals/db/db"
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
	defer rbConn.Close()

	rb := rabbitmq.NewRabbitMQ(rbConn.Conn, rbConn.Ch)

	svc := service.NewService(dbConn, rb)
	svc.StartConsumer()
	svc.PriceCheck("product_check")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go svc.SendForRecheck(ctx, "product_check")

	h := api.NewHandler(log, *svc)
	h.RegisterRoutes(r)

	if err := r.Run(":8000"); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}

}
