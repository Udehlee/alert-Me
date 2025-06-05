package api

import (
	"net/http"

	"github.com/Udehlee/alert-Me/db/db"
	"github.com/Udehlee/alert-Me/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	db  db.Conn
	log zerolog.Logger
}

func NewHandler(db db.Conn, log zerolog.Logger) *Handler {
	return &Handler{
		db:  db,
		log: log,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// r.POST("/submit", h.SubmitProduct)
	r.GET("/signup", h.Signup)
	r.GET("/login", h.Login)

}

func (h *Handler) SubmitProduct(c *gin.Context) {
	var reqProduct models.UrlRequest

	err := c.ShouldBindJSON(&reqProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind URL request"})
		return
	}

}

func (h *Handler) Signup(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unimplemented"})
}

func (h *Handler) Login(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "unimplemented"})
}
