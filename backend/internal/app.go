package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"CorrectQuiz.com/quiz/internal/collection"
	"CorrectQuiz.com/quiz/internal/controller"
	"CorrectQuiz.com/quiz/internal/entity"
	"CorrectQuiz.com/quiz/internal/middleware"
	"CorrectQuiz.com/quiz/internal/service"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"

	firebase "firebase.google.com/go/v4"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const SESSION_SECRET_KEY = "Hfv_${:Fc=-e3Q&Mt5HlXQ{N'ae2M6"

type App struct {
	httpServer  *fiber.App
	database    *gorm.DB
	firebaseApp *firebase.App

	quizService *service.QuizService
	netService  *service.NetService
}

func (a *App) setUpFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	a.firebaseApp = app
}

func (a *App) Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Assuming environment variables are set.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	a.setUpDb()
	a.setUpFirebase()
	a.setUpHttp()

	log.Fatal(a.httpServer.Listen(":" + port))
}

func (a *App) setUpHttp() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://correct-quiz-project.vercel.app, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept,Authorization",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	store := session.New(session.Config{
		CookieName:     "quiz_session",
		CookieSecure:   true,
		CookieSameSite: "None",
		CookieHTTPOnly: true,
		CookieDomain:   "",
	})
	middleware.Store = store

	quizRepo := collection.NewQuizRepository(a.database)
	userRepo := collection.NewGormUserRepository(a.database)
	tokenRepo := collection.NewGormTokenRepository(a.database)
	emailService := service.NewBrevoEmailService()

	firebaseAuthClient, err := a.firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Firebase Auth client: %v\n", err)
	}

	authService := service.NewAuthService(userRepo, firebaseAuthClient, tokenRepo, emailService)
	a.quizService = service.Quiz(quizRepo)
	a.netService = service.Net(a.quizService)

	authController := controller.NewAuthController(authService, store, firebaseAuthClient, userRepo, tokenRepo, emailService)
	gameController := controller.NewGameController(a.netService)
	quizController := controller.Quiz(a.quizService)
	wsController := controller.Ws(a.netService)
	app.Static("/uploads", "./public/uploads")

	app.Post("/api/uploads", func(c *fiber.Ctx) error {
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image upload failed: No file received"})
		}
		uniqueFileName := uuid.New().String() + filepath.Ext(file.Filename)
		url := fmt.Sprintf("http://%s/uploads/%s", c.Hostname(), uniqueFileName)
		return c.JSON(fiber.Map{"url": url})
	})

	app.Post("/api/game/check-name", wsController.CheckPlayerName)
	app.Post("/api/auth/login", authController.Login)
	app.Post("/api/auth/register", authController.Register)
	app.Post("/api/auth/verify-email", authController.VerifyEmailToken)
	app.Post("/api/auth/resend-verification", authController.ResendVerificationEmail)
	app.Post("/api/auth/guest-login", authController.GuestLogin)
	app.Post("/set-initial-claims", authController.HandleSetInitialClaims)
	app.Get("/api/games/:gameCode/export/csv", gameController.ExportGameResultsCSV)
	app.Get("/api/quizzes/:quizId", quizController.GetQuizById)
	app.Get("/ws", websocket.New(wsController.Ws))

	app.Get("/api/users/email/:username", authController.GetUserEmailByUsername)
	app.Post("/api/game/check", wsController.CheckGamePin)

	api := app.Group("/api", middleware.Protected())

	api.Get("/quizzes", quizController.GetCorrect)
	api.Post("/quizzes", quizController.CreateQuiz)
	api.Put("/quizzes/:quizId", quizController.UpdateQuizById)
	api.Delete("/quizzes/:id", quizController.DeleteQuizById)
	api.Delete("/questions/:id", quizController.DeleteQuestionById)

	a.httpServer = app
}

func (a *App) setUpDb() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=Daddymai555 dbname=correct_quiz port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(
		&entity.User{},
		&entity.Quiz{},
		&entity.QuizQuestion{},
		&entity.QuizChoice{},
		&entity.EmailVerificationToken{},
	)

	a.database = db
}
