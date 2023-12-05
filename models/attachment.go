package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Attachment struct {
	ID         uint                  `gorm:"primaryKey;uniqueIndex:udx_attachments"`
	Path       string                `json:"path" gorm:"type:text"`
	DocumentID int                   `json:"document_id"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
	DeletedAt  soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_attachments"`
}
