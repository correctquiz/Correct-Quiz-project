package entity

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name      string         `json:"name"`
	UserID    uint64         `json:"-" gorm:"not null;type:bigint;column:user_id"`
	Questions []QuizQuestion `json:"questions" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"`
}

type QuizQuestion struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name     string       `json:"name"`
	Time     int          `json:"time"`
	ImageUrl string       `json:"imageUrl"`
	Choices  []QuizChoice `json:"choices" gorm:"foreignKey:QuestionID"`
	QuizID   uint
}

type QuizChoice struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name       string  `json:"name"`
	Correct    bool    `json:"correct" gorm:"column:is_correct"`
	ImageUrl   *string `json:"imageUrl,omitempty"`
	QuestionID uint
}
