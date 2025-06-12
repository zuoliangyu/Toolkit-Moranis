package models

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	FilePath  string         `gorm:"size:255;not null" json:"file_path"`
	FolderID  *uint          `gorm:"index" json:"folder_id,omitempty"`            // 文件夹ID，可空外键，添加索引
	Folder    *Folder        `gorm:"foreignKey:FolderID" json:"folder,omitempty"` // 关联的文件夹
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
