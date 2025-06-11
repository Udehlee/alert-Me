package api

import (
	"encoding/json"
	"net/http"

	"github.com/Udehlee/alert-Me/models"
	"github.com/Udehlee/alert-Me/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	rabbit *rabbitmq.RabbitMQ
	log    zerolog.Logger
}

func NewHandler(log zerolog.Logger, rabbit *rabbitmq.RabbitMQ) *Handler {
	return &Handler{
		rabbit: rabbit,
		log:    log,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/submit", h.SubmitProduct)
	r.GET("/signup", h.Signup)
	r.GET("/login", h.Login)

}

// SubmitProduct handles product url request
func (h *Handler) SubmitProduct(c *gin.Context) {
	var reqProduct models.UrlRequest

	if err := c.ShouldBindJSON(&reqProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind URL request"})
		return
	}

	body, err := json.Marshal(reqProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal reqProduct request"})
		return
	}

	err = h.rabbit.PublishToQueue("product_url_queue", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish URL to queue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL received and processing started"})
}

func (h *Handler) Signup(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unimplemented"})
}

func (h *Handler) Login(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unimplemented"})
}
