package collection

import (
	"context"

	"CorrectQuiz.com/quiz/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByFirebaseUID(ctx context.Context, uid string) (*entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	UpdateUser(user *entity.User) error
	GetUserByUsername(username string) (*entity.User, error)
	IsUsernameTaken(username string) (bool, error)
	FindByID(userID uint) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	IsEmailTaken(email string) (bool, error)
}

type UserGormRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(database *gorm.DB) UserRepository {
	return &UserGormRepository{
		db: database,
	}
}

func (r *UserGormRepository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserGormRepository) GetUserByFirebaseUID(ctx context.Context, uid string) (*entity.User, error) {
	var user entity.User
	result := r.db.WithContext(ctx).Where("firebase_uid = ?", uid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserGormRepository) GetUserByID(id uint) (*entity.User, error) {
	var user entity.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserGormRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserGormRepository) UpdateUser(user *entity.User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r *UserGormRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserGormRepository) IsUsernameTaken(username string) (bool, error) {
	var count int64

	result := r.db.Model(&entity.User{}).Where("username = ? AND deleted_at IS NULL", username).Count(&count)
	if result.Error != nil {

		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return count > 0, nil
}

func (r *UserGormRepository) IsEmailTaken(email string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserGormRepository) FindByID(userID uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, userID).Error
	return &user, err
}
