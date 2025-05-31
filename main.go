package main

import (
	"os"

	"github.com/Udehlee/alert-Me/api"
	"github.com/Udehlee/alert-Me/client"
	"github.com/Udehlee/alert-Me/db/db"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	r := gin.Default()
	log := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	Client := client.NewEbayClient(os.Getenv("ACCESS_TOKEN"))

	db, err := db.InitConnectDB()
	if err != nil {
		log.Err(err).Msg("failed to connect to database")

	}

	h := api.NewHandler(db, Client, log)
	h.RegisterRoutes(r)

}
