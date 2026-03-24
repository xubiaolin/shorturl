package handlers

import (
	"net/http"

	"shorturl/service"

	"github.com/gin-gonic/gin"
)

type RedirectHandler struct {
	service *service.ShortURLService
}

func NewRedirectHandler() *RedirectHandler {
	return &RedirectHandler{
		service: service.NewShortURLService(),
	}
}

// Redirect 重定向到原始 URL
func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("code")

	url, err := h.service.GetByShortCode(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}
