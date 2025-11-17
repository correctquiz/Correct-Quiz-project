package collection

import (
	"CorrectQuiz.com/quiz/internal/entity"
	"gorm.io/gorm"
)

type QuizRepository interface {
	InsertQuiz(quiz *entity.Quiz) error
	GetCorrect(userID uint) ([]entity.Quiz, error)
	GetQuizById(id uint) (*entity.Quiz, error)
	UpdateQuiz(quiz entity.Quiz) error
	DeleteQuestionById(id uint) error
	DeleteQuizById(id uint) error
	GetQuizzesByUserID(userID uint) ([]entity.Quiz, error)
}

type quizGormRepository struct {
	db *gorm.DB
}

func NewQuizRepository(database *gorm.DB) QuizRepository {
	return &quizGormRepository{
		db: database,
	}
}

func (r *quizGormRepository) GetQuizzesByUserID(userID uint) ([]entity.Quiz, error) {
	var quizzes []entity.Quiz

	result := r.db.
		Preload("Questions.Choices").
		Where("user_id = ?", userID).
		Find(&quizzes)

	if result.Error != nil {
		return nil, result.Error
	}

	return quizzes, nil
}

func (r *quizGormRepository) DeleteQuestionById(id uint) error {
	return r.db.Delete(&entity.QuizQuestion{}, id).Error
}

func (r *quizGormRepository) DeleteQuizById(id uint) error {
	return r.db.Delete(&entity.Quiz{}, id).Error
}

func (r *quizGormRepository) InsertQuiz(quiz *entity.Quiz) error {
	result := r.db.Create(quiz)
	return result.Error
}

func (r *quizGormRepository) GetCorrect(userID uint) ([]entity.Quiz, error) {
	var quizzes []entity.Quiz
	result := r.db.
		Preload("Questions.Choices").
		Preload("Questions").
		Where("user_id = ?", userID).
		Find(&quizzes)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return quizzes, nil
}

func (r *quizGormRepository) GetQuizById(id uint) (*entity.Quiz, error) {
	var quiz entity.Quiz
	result := r.db.Preload("Questions.Choices").Preload("Questions").First(&quiz, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &quiz, nil
}

func (r *quizGormRepository) UpdateQuiz(quiz entity.Quiz) error {
	return r.db.
		Omit("CreatedAt").
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&quiz).Error
}
