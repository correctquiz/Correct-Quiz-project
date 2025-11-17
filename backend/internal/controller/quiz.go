package controller

import (
	"strconv"

	"CorrectQuiz.com/quiz/internal/entity"
	"CorrectQuiz.com/quiz/internal/service"
	"github.com/gofiber/fiber/v2"
)

type QuizController struct {
	quizService *service.QuizService
}

func Quiz(svc *service.QuizService) *QuizController {
	return &QuizController{quizService: svc}
}

func (c *QuizController) GetQuizById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("quizId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format. ID must be a number.",
		})
	}

	quiz, err := c.quizService.GetQuizById(uint(id))

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Quiz not found"})
	}

	return ctx.JSON(quiz)
}

type UpdateQuizRequest struct {
	Name      string                `json:"name"`
	Questions []entity.QuizQuestion `json:"questions"`
}

type CreateQuizRequest struct {
	Name      string                `json:"name"`
	Questions []entity.QuizQuestion `json:"questions"`
}

func (c *QuizController) UpdateQuizById(ctx *fiber.Ctx) error {
	quizIdStr := ctx.Params("quizId")

	quizId, err := strconv.ParseUint(quizIdStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format. ID must be a number.",
		})
	}

	userIDRaw := ctx.Locals("user_id")

	userID, ok := userIDRaw.(uint)

	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in session"})
	}

	var req UpdateQuizRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body.",
		})
	}

	if err := c.quizService.UpdateQuiz(uint(quizId), uint64(userID), req.Name, req.Questions); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *QuizController) CreateQuiz(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(uint)
	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found or invalid in session"})
	}

	var req CreateQuizRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse body"})
	}
	newQuiz := entity.Quiz{
		Name:      req.Name,
		Questions: req.Questions,
		UserID:    uint64(userID),
	}

	createdQuiz, err := c.quizService.CreateQuiz(newQuiz)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdQuiz)
}

func (c *QuizController) DeleteQuestionById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if err := c.quizService.DeleteQuestionById(uint(id)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *QuizController) GetCorrect(ctx *fiber.Ctx) error {
	userIDRaw := ctx.Locals("user_id")
	userID, ok := userIDRaw.(uint)
	if !ok || userID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in session"})
	}
	quizzes, err := c.quizService.GetQuizzesByUserID(userID)
	if err != nil {
		return err
	}

	return ctx.JSON(quizzes)
}

func (c *QuizController) DeleteQuizById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := c.quizService.DeleteQuizById(uint(id)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
