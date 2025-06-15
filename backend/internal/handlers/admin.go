package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"servermanager/internal/config"
	"servermanager/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// 管理员登录
func (h *AdminHandler) Login(c *gin.Context) {
	var loginData struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 这里应该使用更安全的方式验证密码
	if loginData.Password != config.Config.Server.Admin.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成JWT token
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// 生成JWT token
func generateToken() (string, error) {
	claims := jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.Server.Admin.JWTSecret))
}

// 上传资料
func (h *AdminHandler) UploadMaterial(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "资料名称不能为空"})
		return
	}

	// 确保上传目录存在
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	baseName := filepath.Base(file.Filename[:len(file.Filename)-len(ext)])
	filename := fmt.Sprintf("%s_%d%s", baseName, time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 保存到数据库
	material := models.Material{
		Name:     name,
		FilePath: filename,
	}

	if err := h.db.Create(&material).Error; err != nil {
		os.Remove(filepath) // 删除已上传的文件
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存资料信息失败"})
		return
	}

	c.JSON(http.StatusOK, material)
}

// 删除资料
func (h *AdminHandler) DeleteMaterial(c *gin.Context) {
	id := c.Param("id")
	var material models.Material
	if err := h.db.First(&material, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}

	// 删除文件
	filepath := filepath.Join("uploads", material.FilePath)
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}

	// 删除数据库记录
	if err := h.db.Delete(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除资料记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 获取文件夹树形结构
func (h *AdminHandler) GetFolders(c *gin.Context) {
	var folders []models.Folder
	if err := h.db.Preload("Children").Where("parent_id IS NULL").Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件夹列表失败"})
		return
	}
	c.JSON(http.StatusOK, folders)
}

// 创建文件夹
func (h *AdminHandler) CreateFolder(c *gin.Context) {
	var folderData struct {
		Name        string `json:"name" binding:"required"`
		ParentID    *uint  `json:"parent_id"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&folderData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证文件夹名称
	if strings.TrimSpace(folderData.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件夹名称不能为空"})
		return
	}

	folder := models.Folder{
		Name:        strings.TrimSpace(folderData.Name),
		ParentID:    folderData.ParentID,
		Description: folderData.Description,
		Level:       0,
		FolderPath:  folderData.Name,
	}

	// 如果有父文件夹，验证层级和计算路径
	if folderData.ParentID != nil {
		var parentFolder models.Folder
		if err := h.db.First(&parentFolder, *folderData.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "父文件夹不存在"})
			return
		}

		// 检查层级限制
		if parentFolder.Level >= 99 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件夹层级不能超过99层"})
			return
		}

		folder.Level = parentFolder.Level + 1
		folder.FolderPath = parentFolder.FolderPath + "/" + folder.Name
	}

	// 检查同级文件夹名称是否重复
	var existingFolder models.Folder
	query := h.db.Where("name = ?", folder.Name)
	if folder.ParentID != nil {
		query = query.Where("parent_id = ?", *folder.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if err := query.First(&existingFolder).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "同级目录下已存在同名文件夹"})
		return
	}

	if err := h.db.Create(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件夹失败"})
		return
	}

	c.JSON(http.StatusOK, folder)
}

// 重命名文件夹
func (h *AdminHandler) UpdateFolder(c *gin.Context) {
	id := c.Param("id")
	folderID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件夹ID"})
		return
	}

	var updateData struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 验证文件夹名称
	newName := strings.TrimSpace(updateData.Name)
	if newName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件夹名称不能为空"})
		return
	}

	var folder models.Folder
	if err := h.db.First(&folder, uint(folderID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件夹不存在"})
		return
	}

	// 检查同级文件夹名称是否重复
	if folder.Name != newName {
		var existingFolder models.Folder
		query := h.db.Where("name = ? AND id != ?", newName, folder.ID)
		if folder.ParentID != nil {
			query = query.Where("parent_id = ?", *folder.ParentID)
		} else {
			query = query.Where("parent_id IS NULL")
		}

		if err := query.First(&existingFolder).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "同级目录下已存在同名文件夹"})
			return
		}
	}

	oldPath := folder.FolderPath
	oldName := folder.Name

	// 更新文件夹信息
	folder.Name = newName
	folder.Description = updateData.Description

	// 重新计算路径
	if folder.ParentID != nil {
		var parentFolder models.Folder
		if err := h.db.First(&parentFolder, *folder.ParentID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "父文件夹不存在"})
			return
		}
		folder.FolderPath = parentFolder.FolderPath + "/" + newName
	} else {
		folder.FolderPath = newName
	}

	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新当前文件夹
	if err := tx.Save(&folder).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文件夹失败"})
		return
	}

	// 如果名称发生变化，需要更新所有子文件夹的路径
	if oldName != newName {
		if err := h.updateChildrenPaths(tx, folder.ID, oldPath, folder.FolderPath); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新子文件夹路径失败"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, folder)
}

// 删除文件夹
func (h *AdminHandler) DeleteFolder(c *gin.Context) {
	id := c.Param("id")
	folderID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件夹ID"})
		return
	}

	var folder models.Folder
	if err := h.db.First(&folder, uint(folderID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件夹不存在"})
		return
	}

	// 检查是否有子文件夹
	var childCount int64
	if err := h.db.Model(&models.Folder{}).Where("parent_id = ?", folder.ID).Count(&childCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查子文件夹失败"})
		return
	}

	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件夹不为空，请先删除子文件夹"})
		return
	}

	// 检查是否有资料文件
	var materialCount int64
	if err := h.db.Model(&models.Material{}).Where("folder_id = ?", folder.ID).Count(&materialCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查资料文件失败"})
		return
	}

	if materialCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件夹不为空，请先删除或移动其中的资料文件"})
		return
	}

	// 删除文件夹
	if err := h.db.Delete(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件夹失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 辅助函数：递归更新子文件夹路径
func (h *AdminHandler) updateChildrenPaths(tx *gorm.DB, parentID uint, oldParentPath, newParentPath string) error {
	var children []models.Folder
	if err := tx.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		// 更新子文件夹路径
		newChildPath := strings.Replace(child.FolderPath, oldParentPath, newParentPath, 1)
		if err := tx.Model(&child).Update("folder_path", newChildPath).Error; err != nil {
			return err
		}

		// 递归更新子文件夹的子文件夹
		if err := h.updateChildrenPaths(tx, child.ID, child.FolderPath, newChildPath); err != nil {
			return err
		}
	}

	return nil
}
