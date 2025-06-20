package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Udehlee/alert-Me/models"
	"github.com/Udehlee/alert-Me/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	service service.Service
	log     zerolog.Logger
}

func NewHandler(log zerolog.Logger, svc service.Service) *Handler {
	return &Handler{
		service: svc,
		log:     log,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/submit", h.SubmitProduct)
	r.GET("/", h.Index)

}

func (h *Handler) Index(c *gin.Context) {
	c.String(200, "Welcome Home, my gee")
}

// SubmitProduct handles product url request
func (h *Handler) SubmitProduct(c *gin.Context) {
	var reqProduct models.UrlRequest

	if err := c.ShouldBindJSON(&reqProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind product_url request"})
		return
	}

	body, err := json.Marshal(reqProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal reqProduct request"})
		return
	}

	err = h.service.Rabbit.PublishToQueue("product_url_queue", body)
	if err != nil {
		// log.Printf(" Failed to publish message to queue: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue the URL"})
		return
	}

	log.Printf("Successfully published message for URL: %s", reqProduct.URL)
	c.JSON(http.StatusOK, gin.H{"message": "URL received and processing started"})
}
