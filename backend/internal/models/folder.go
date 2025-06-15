package models

import (
	"time"
	"gorm.io/gorm"
)

type Folder struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`                    // 文件夹名称
	ParentID    *uint          `gorm:"index" json:"parent_id,omitempty"`                 // 父文件夹ID，可空外键，添加索引
	Parent      *Folder        `gorm:"foreignKey:ParentID" json:"parent,omitempty"`      // 关联的父文件夹
	Children    []Folder       `gorm:"foreignKey:ParentID" json:"children,omitempty"`    // 关联的子文件夹
	FolderPath  string         `gorm:"size:500;index" json:"folder_path"`                // 完整路径，用于快速查询
	Level       int            `gorm:"default:0;check:level <= 99" json:"level"`         // 文件夹层级，限制最大99层
	Description string         `gorm:"size:500" json:"description"`                      // 文件夹描述
	IsDefault   bool           `gorm:"default:false" json:"is_default"`                  // 是否为默认文件夹
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                                   // 软删除
	Materials   []Material     `gorm:"foreignKey:FolderID" json:"materials,omitempty"`   // 关联的资料
}
