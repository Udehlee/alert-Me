package api

import (
	"net/http"

	"github.com/Udehlee/alert-Me/client"
	"github.com/Udehlee/alert-Me/db/db"
	"github.com/Udehlee/alert-Me/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	db     db.Conn
	client *client.EbayClient
	log    zerolog.Logger
}

func NewHandler(db db.Conn, ec *client.EbayClient, log zerolog.Logger) *Handler {
	return &Handler{
		db:     db,
		client: ec,
		log:    log,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/search", h.SearchProduct)
	r.GET("/signup", h.Signup)
	r.GET("/login", h.Login)

}

func (h *Handler) SearchProduct(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		h.log.Warn().Msg("Missing query parameter 'q'")
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	product, err := h.client.GetProduct(query)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get product from eBay")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	selectdProduct := models.SelectedProduct{
		ProductID:    product.ItemID,
		ProductName:  product.Title,
		CurrentPrice: product.Price.Value,
	}

	err = h.db.SelectedProduct(selectdProduct)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save selected product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully gotten the product details"})
}

func (h *Handler) Signup(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unimplemented"})
}

func (h *Handler) Login(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unimplemented"})
}
