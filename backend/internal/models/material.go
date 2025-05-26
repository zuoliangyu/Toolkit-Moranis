package models

import (
	"time"
	"gorm.io/gorm"
)

type Material struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	FilePath  string    `gorm:"size:255;not null" json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
} 