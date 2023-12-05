package models

import (
	"time"

	"github.com/mhosseintaher/kit/dtp"
	"gorm.io/plugin/soft_delete"
)

type Category struct {
	ID        uint                  `gorm:"primaryKey"`
	Children  []*Category           `json:"children" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Parent    *Category             `json:"parent" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ParentID  dtp.NullInt64         `json:"parent_id"`
	Title     string                `json:"title" gorm:"size:100;uniqueIndex:udx_categories;"`
	User      *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID    int                   `json:"user_id"`
	Rate      float64               `json:"rate"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_categories"`
}
