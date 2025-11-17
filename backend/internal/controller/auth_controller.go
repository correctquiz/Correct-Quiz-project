package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	firebaseAuth "firebase.google.com/go/v4/auth"

	"CorrectQuiz.com/quiz/internal/collection"
	"CorrectQuiz.com/quiz/internal/entity"
	"CorrectQuiz.com/quiz/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
)

type AuthController struct {
	authService  service.AuthServiceInterface
	sessionStore *session.Store
	authClient   *firebaseAuth.Client
	userRepo     collection.UserRepository
	tokenRepo    collection.TokenRepository
	emailService service.EmailServiceInterface
}

func NewAuthController(
	svc service.AuthServiceInterface,
	store *session.Store,
	ac *firebaseAuth.Client,
	ur collection.UserRepository,
	tr collection.TokenRepository,
	es service.EmailServiceInterface,
) *AuthController {
	return &AuthController{
		authService:  svc,
		sessionStore: store,
		authClient:   ac,
		userRepo:     ur,
		tokenRepo:    tr,
		emailService: es,
	}
}

type SetInitialClaimsRequest struct {
	FirebaseUID string `json:"firebaseUid"`
	Username    string `json:"username"`
	Role        string `json:"role"`
}

type LoginRequest struct {
	IdToken string `json:"idToken"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirebaseUID string `json:"firebaseUid"`
	Role        string `json:"role"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type ResendRequest struct {
	Email string `json:"email"`
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	isTaken, err := c.userRepo.IsUsernameTaken(req.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error checking username"})
	}
	if isTaken {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username นี้ ใช้ไปแล้ว"})
	}

	isEmailTaken, err := c.userRepo.IsEmailTaken(req.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error checking email"})
	}
	if isEmailTaken {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already taken"})
	}

	claims := map[string]interface{}{
		"username": req.Username,
		"role":     req.Role,
	}
	err = c.authClient.SetCustomUserClaims(context.Background(), req.FirebaseUID, claims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set user claims"})
	}

	user := entity.User{
		Email:         req.Email,
		Role:          req.Role,
		Username:      req.Username,
		FirebaseUID:   req.FirebaseUID,
		EmailVerified: false,
	}
	createdUser, err := c.authService.Register(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tokenString := uuid.New().String()
	verificationToken := &entity.EmailVerificationToken{
		UserID:    createdUser.ID,
		Token:     tokenString,
		CreatedAt: time.Now(),
	}
	if err := c.tokenRepo.Create(verificationToken); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create verification token"})
	}

	verificationLink := fmt.Sprintf("http://localhost:5173/#/verify-email?token=%s", tokenString)
	htmlContent := fmt.Sprintf(
		`<strong>ยินดีต้อนรับสู่ Correct Quiz!</strong>
         <p>กรุณาคลิกลิงก์ด้านล่างเพื่อยืนยันอีเมลของคุณ (ลิงก์นี้จะหมดอายุใน 5 นาที):</p>
         <a href="%s" target="_blank">คลิกที่นี่เพื่อยืนยัน</a>`,
		verificationLink,
	)

	err = c.emailService.SendEmail(
		createdUser.Email,
		createdUser.Username,
		"ยืนยันอีเมลของคุณ - Correct Quiz (หมดอายุใน 5 นาที)",
		htmlContent,
	)
	if err != nil {
		log.Printf("CRITICAL: Failed to send verification email to %s: %v", createdUser.Email, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdUser)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {

	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}

	user, err := c.authService.Login(ctx.Context(), req.IdToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials or token"})
	}

	sess, err := c.sessionStore.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session error"})
	}

	sess.Set("user_id", user.ID)
	if err := sess.Save(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session save error"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   "session_created",
	})
}

func (c *AuthController) ResendVerificationEmail(ctx *fiber.Ctx) error {
	var req ResendRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing email"})
	}

	err := c.authService.ResendVerification(req.Email)
	if err != nil {
		// ส่ง Error ที่ service สร้างไว้ (เช่น "อีเมลนี้ถูกยืนยันตัวตนไปแล้ว")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Verification email resent successfully"})
}

func (c *AuthController) HandleSetInitialClaims(ctx *fiber.Ctx) error {
	var req SetInitialClaimsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.FirebaseUID == "" || req.Username == "" || req.Role == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	isTaken, err := c.userRepo.IsUsernameTaken(req.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error checking username"})
	}
	if isTaken {

		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username นี้ ใช้ไปแล้ว"})
	}

	claims := map[string]interface{}{
		"username": req.Username,
		"role":     req.Role,
	}

	err = c.authClient.SetCustomUserClaims(context.Background(), req.FirebaseUID, claims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set user claims"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *AuthController) GetUserEmailByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	user, err := c.authService.GetUserByUsername(username)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(fiber.Map{"email": user.Email})
}

func (ctl *AuthController) VerifyEmailToken(c *fiber.Ctx) error {
	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	tokenData, err := ctl.tokenRepo.FindByToken(req.Token)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Token ไม่ถูกต้อง หรือถูกใช้ไปแล้ว"})
	}

	expirationTime := 5 * time.Minute
	if time.Since(tokenData.CreatedAt) > expirationTime {
		ctl.tokenRepo.Delete(tokenData.ID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token หมดอายุแล้ว กรุณาสมัครสมาชิกใหม่อีกครั้ง"})
	}

	err = ctl.authService.VerifyUserEmailInFirebase(tokenData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify user in Firebase"})
	}

	ctl.tokenRepo.Delete(tokenData.ID)

	return c.JSON(fiber.Map{"message": "Email verified successfully"})
}
