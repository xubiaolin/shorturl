package handlers

import (
	"net/http"
	"strconv"
	"time"

	"shorturl/service"

	"github.com/gin-gonic/gin"
)

type ShortURLHandler struct {
	service *service.ShortURLService
}

func NewShortURLHandler() *ShortURLHandler {
	return &ShortURLHandler{
		service: service.NewShortURLService(),
	}
}

// CreateShortURL 创建短链
// @Summary 创建短链
// @Tags 短链管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body CreateShortURLRequest true "短链信息"
// @Success 200 {object} Response
// @Router /api/v1/shorturls [post]
func (h *ShortURLHandler) CreateShortURL(c *gin.Context) {
	var req struct {
		OriginalURL string `json:"original_url" binding:"required"`
		CustomCode  string `json:"custom_code"`
		ExpiresAt   *int64 `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t := time.Unix(*req.ExpiresAt, 0)
		expiresAt = &t
	}

	url, err := h.service.CreateShortURL(req.OriginalURL, req.CustomCode, expiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    url,
	})
}

// UpdateShortURL 更新短链
// @Summary 更新短链
// @Tags 短链管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "短链 ID"
// @Param body body UpdateShortURLRequest true "短链信息"
// @Success 200 {object} Response
// @Router /api/v1/shorturls/:id [put]
func (h *ShortURLHandler) UpdateShortURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		OriginalURL string `json:"original_url"`
		CustomCode  string `json:"custom_code"`
		ExpiresAt   *int64 `json:"expires_at"`
		IsActive    *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
			"data":    nil,
		})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t := time.Unix(*req.ExpiresAt, 0)
		expiresAt = &t
	}

	url, err := h.service.UpdateShortURL(uint(id), req.OriginalURL, req.CustomCode, expiresAt, req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    url,
	})
}

// DeleteShortURL 删除短链
// @Summary 删除短链
// @Tags 短链管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "短链 ID"
// @Success 200 {object} Response
// @Router /api/v1/shorturls/:id [delete]
func (h *ShortURLHandler) DeleteShortURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
			"data":    nil,
		})
		return
	}

	if err := h.service.DeleteShortURL(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}

// GetShortURL 获取短链详情
// @Summary 获取短链详情
// @Tags 短链管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "短链 ID"
// @Success 200 {object} Response
// @Router /api/v1/shorturls/:id [get]
func (h *ShortURLHandler) GetShortURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
			"data":    nil,
		})
		return
	}

	url, err := h.service.GetShortURL(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    url,
	})
}

// ListShortURLs 获取短链列表
// @Summary 获取短链列表
// @Tags 短链管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} Response
// @Router /api/v1/shorturls [get]
func (h *ShortURLHandler) ListShortURLs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	urls, total, err := h.service.ListShortURLs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":  urls,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// StatsShortURL 获取短链统计
// @Summary 获取短链统计
// @Tags 短链管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "短链 ID"
// @Success 200 {object} Response
// @Router /api/v1/shorturls/:id/stats [get]
func (h *ShortURLHandler) StatsShortURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
			"data":    nil,
		})
		return
	}

	stats, err := h.service.StatsShortURL(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    stats,
	})
}
