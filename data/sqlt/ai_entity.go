package sqlt

import (
	"gorm.io/gorm"
)

type DetectedObject struct {
	Score float32 `json:"score"`
	Label string  `json:"label" gorm:"index:idx_query"` //Index
}

type AiEntity struct {
	gorm.Model
	BaseEntity
	DetectedObject
}
