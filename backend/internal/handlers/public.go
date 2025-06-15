package handlers

import (
	"net/http"
	"path/filepath"
	"servermanager/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PublicHandler struct {
	db *gorm.DB
}

func NewPublicHandler(db *gorm.DB) *PublicHandler {
	return &PublicHandler{db: db}
}

// 获取所有资料列表，支持按文件夹筛选
func (h *PublicHandler) GetMaterials(c *gin.Context) {
	var materials []models.Material
	query := h.db.Preload("Folder")

	// 支持按文件夹ID筛选
	folderIDStr := c.Query("folder_id")
	if folderIDStr != "" {
		folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件夹ID"})
			return
		}
		query = query.Where("folder_id = ?", uint(folderID))
	}

	if err := query.Find(&materials).Error; err != nil {
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

// 获取文件夹树形结构
func (h *PublicHandler) GetFolders(c *gin.Context) {
	var folders []models.Folder
	if err := h.db.Preload("Children").Where("parent_id IS NULL").Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件夹列表失败"})
		return
	}
	c.JSON(http.StatusOK, folders)
}

// 获取指定文件夹下的资料
func (h *PublicHandler) GetFolderMaterials(c *gin.Context) {
	id := c.Param("id")
	folderID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件夹ID"})
		return
	}

	// 验证文件夹是否存在
	var folder models.Folder
	if err := h.db.First(&folder, uint(folderID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件夹不存在"})
		return
	}

	// 获取文件夹下的资料
	var materials []models.Material
	if err := h.db.Preload("Folder").Where("folder_id = ?", uint(folderID)).Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取资料列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folder":    folder,
		"materials": materials,
	})
}
