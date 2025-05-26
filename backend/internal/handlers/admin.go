package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"servermanager/internal/config"
	"servermanager/internal/models"
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
