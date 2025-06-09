package main

import (
	"os"

	"github.com/Udehlee/alert-Me/api"
	"github.com/Udehlee/alert-Me/db/db"
	"github.com/Udehlee/alert-Me/pkg/rabbitMQ"
	"github.com/Udehlee/alert-Me/pkg/utils"
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

	rbConn, err := rabbitMQ.ConnectRabbitMQ()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to RabbitMQ")
	}

	rb := rabbitMQ.NewRabbitMQ(rbConn.Conn, rbConn.Ch, dbConn)

	go rb.Consumer("product-url", utils.ExtractProduct)

	h := api.NewHandler(log, &rbConn)
	h.RegisterRoutes(r)

	if err := r.Run(":8000"); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}

}
