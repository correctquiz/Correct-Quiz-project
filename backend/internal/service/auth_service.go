package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"CorrectQuiz.com/quiz/internal/collection"
	"CorrectQuiz.com/quiz/internal/entity"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	Login(ctx context.Context, idToken string) (*entity.User, error)
	VerifyUserEmailInFirebase(userID uint) error
	ResendVerification(email string) error
}

type AuthService struct {
	userRepo     collection.UserRepository
	authClient   *firebaseAuth.Client
	tokenRepo    collection.TokenRepository
	emailService EmailServiceInterface
}

func NewAuthService(ur collection.UserRepository, ac *firebaseAuth.Client, tr collection.TokenRepository, es EmailServiceInterface) AuthServiceInterface {
	return &AuthService{
		userRepo:     ur,
		authClient:   ac,
		tokenRepo:    tr,
		emailService: es,
	}
}

func (s *AuthService) ResendVerification(email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("ไม่พบผู้ใช้ที่ใช้อีเมลนี้")
	}

	if user.EmailVerified {
		return fmt.Errorf("อีเมลนี้ถูกยืนยันตัวตนไปแล้ว")
	}

	if err := s.tokenRepo.DeleteByUserID(user.ID); err != nil {
		log.Printf("Warning: failed to delete old tokens for user %d: %v", user.ID, err)
	}

	tokenString := uuid.New().String()
	verificationToken := &entity.EmailVerificationToken{
		UserID:    user.ID,
		Token:     tokenString,
		CreatedAt: time.Now(),
	}
	if err := s.tokenRepo.Create(verificationToken); err != nil {
		return fmt.Errorf("failed to create new verification token")
	}

	verificationLink := fmt.Sprintf("http://localhost:5173/#/verify-email?token=%s", tokenString)
	htmlContent := fmt.Sprintf(
		`<strong>คุณขอส่งอีเมลยืนยันอีกครั้ง</strong>
         <p>กรุณาคลิกลิงก์ด้านล่างเพื่อยืนยันอีเมล (ลิงก์นี้จะหมดอายุใน 5 นาที):</p>
         <a href="%s" target="_blank">คลิกที่นี่เพื่อยืนยัน</a>`,
		verificationLink,
	)

	err = s.emailService.SendEmail(
		user.Email,
		user.Username,
		"ส่งซ้ำ: ยืนยันอีเมลของคุณ - Correct Quiz (หมดอายุใน 5 นาที)",
		htmlContent,
	)
	if err != nil {
		log.Printf("CRITICAL: Failed to resend verification email to %s: %v", user.Email, err)
		return fmt.Errorf("failed to send email")
	}

	return nil
}

func (s *AuthService) Register(ctx context.Context, user entity.User) (*entity.User, error) {

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *AuthService) VerifyUserEmailInFirebase(userID uint) error {

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	params := (&firebaseAuth.UserToUpdate{}).EmailVerified(true)
	_, err = s.authClient.UpdateUser(context.Background(), user.FirebaseUID, params)

	if err != nil {
		return fmt.Errorf("failed to update user in Firebase: %w", err)
	}

	user.EmailVerified = true

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user in Postgres: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, idToken string) (*entity.User, error) {
	token, err := s.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {

		return nil, errors.New("invalid or expired firebase token")
	}

	uid := token.UID
	claims := token.Claims

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("email not found or invalid in token claims")
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return nil, errors.New("username claim not found, invalid, or empty")
	}
	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return nil, errors.New("role claim not found, invalid, or empty")
	}

	user, err := s.userRepo.GetUserByFirebaseUID(ctx, uid)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			isTaken, checkErr := s.userRepo.IsUsernameTaken(username)
			if checkErr != nil {
				return nil, errors.New("database error during login (username check)")
			}
			if isTaken {
				return nil, errors.New("username conflict detected during first login")
			}

			newUser := entity.User{
				FirebaseUID:   uid,
				Email:         email,
				Username:      username,
				Role:          role,
				EmailVerified: true,
			}

			createdUser, createErr := s.userRepo.CreateUser(ctx, newUser)
			if createErr != nil {
				return nil, errors.New("failed to register user in database during first login")
			}

			return createdUser, nil
		}

		return nil, errors.New("database error during login")
	}
	return user, nil
}

func (s *AuthService) GetUserByUsername(username string) (*entity.User, error) {
	return s.userRepo.GetUserByUsername(username)
}
