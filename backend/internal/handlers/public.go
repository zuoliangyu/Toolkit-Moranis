package handlers

import (
	"net/http"
	"path/filepath"
	"servermanager/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PublicHandler struct {
	db *gorm.DB
}

func NewPublicHandler(db *gorm.DB) *PublicHandler {
	return &PublicHandler{db: db}
}

// 获取所有资料列表
func (h *PublicHandler) GetMaterials(c *gin.Context) {
	var materials []models.Material
	if err := h.db.Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资料列表失败"})
		return
	}
	c.JSON(http.StatusOK, materials)
}

// 下载资料
func (h *PublicHandler) DownloadMaterial(c *gin.Context) {
	id := c.Param("id")
	var material models.Material
	if err := h.db.First(&material, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}

	filePath := material.FilePath
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join("uploads", filePath)
	}

	c.File(filePath)
}
