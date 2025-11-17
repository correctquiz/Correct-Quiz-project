package controller

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"CorrectQuiz.com/quiz/internal/service"
	"github.com/gofiber/fiber/v2"
)

type GameController struct {
	netService *service.NetService
}

func NewGameController(ns *service.NetService) *GameController {
	return &GameController{netService: ns}
}

func (gc *GameController) ExportGameResultsCSV(c *fiber.Ctx) error {
	gameCode := c.Params("gameCode")

	game, err := gc.netService.FindActiveGameByCode(gameCode)
	if err != nil || game == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Active game with this code not found"})
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"Player Name", "Final Score", "CorrectStreak"}
	if err := writer.Write(header); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	players := game.Players

	for _, player := range players {
		correctCount := 0
		if count, ok := game.GetCorrectAnswerCount(player.Id); ok {
			correctCount = count
		}
		row := []string{
			player.Name,
			strconv.Itoa(player.Points),
			strconv.Itoa(correctCount),
		}
		if err := writer.Write(row); err != nil {
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"quiz_results_%s.csv\"", gameCode))
	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")

	return c.Send(buf.Bytes())
}
