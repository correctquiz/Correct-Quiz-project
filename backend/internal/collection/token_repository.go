package collection

import (
	"CorrectQuiz.com/quiz/internal/entity"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(token *entity.EmailVerificationToken) error
	FindByToken(token string) (*entity.EmailVerificationToken, error)
	Delete(tokenID uint) error
	DeleteByUserID(userID uint) error
}

type gormTokenRepository struct {
	db *gorm.DB
}

func NewGormTokenRepository(db *gorm.DB) TokenRepository {
	return &gormTokenRepository{db: db}
}

func (r *gormTokenRepository) Create(token *entity.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

func (r *gormTokenRepository) FindByToken(token string) (*entity.EmailVerificationToken, error) {
	var verificationToken entity.EmailVerificationToken
	err := r.db.Where("token = ?", token).First(&verificationToken).Error
	return &verificationToken, err
}

func (r *gormTokenRepository) Delete(tokenID uint) error {
	return r.db.Delete(&entity.EmailVerificationToken{}, tokenID).Error
}

func (r *gormTokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.EmailVerificationToken{}).Error
}
