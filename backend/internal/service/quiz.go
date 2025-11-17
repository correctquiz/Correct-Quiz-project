package service

import (
	"CorrectQuiz.com/quiz/internal/collection"
	"CorrectQuiz.com/quiz/internal/entity"
	"gorm.io/gorm"
)

type QuizService struct {
	quizCollection collection.QuizRepository
}

type QuizRepository struct {
	DB *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *QuizRepository {
	return &QuizRepository{DB: db}
}

func Quiz(quizRepo collection.QuizRepository) *QuizService {
	return &QuizService{
		quizCollection: quizRepo,
	}
}

func (s *QuizService) GetCorrect(userID uint) ([]entity.Quiz, error) {
	quizzes, err := s.quizCollection.GetCorrect(userID)

	if err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (s *QuizService) GetQuizzesByUserID(userID uint) ([]entity.Quiz, error) {
	quizzes, err := s.quizCollection.GetQuizzesByUserID(userID)
	if err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (s *QuizService) GetQuizById(id uint) (*entity.Quiz, error) {
	return s.quizCollection.GetQuizById(id)
}

func (s *QuizService) UpdateQuiz(id uint, userID uint64, name string, questions []entity.QuizQuestion) error {
	quizToUpdate := entity.Quiz{
		ID:        id,
		Name:      name,
		Questions: questions,
		UserID:    userID,
	}

	for i := range quizToUpdate.Questions {
		quizToUpdate.Questions[i].QuizID = id
	}

	return s.quizCollection.UpdateQuiz(quizToUpdate)
}

func (s *QuizService) DeleteQuestionById(id uint) error {
	return s.quizCollection.DeleteQuestionById(id)
}

func (s *QuizService) CreateQuiz(quiz entity.Quiz) (*entity.Quiz, error) {
	if err := s.quizCollection.InsertQuiz(&quiz); err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (s *QuizService) DeleteQuizById(id uint) error {
	return s.quizCollection.DeleteQuizById(id)
}
