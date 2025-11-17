package controller

import (
	"CorrectQuiz.com/quiz/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebsocketController struct {
	netService *service.NetService
}

type CheckNamePayload struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func Ws(netService *service.NetService) WebsocketController {
	return WebsocketController{
		netService: netService,
	}
}

func (c WebsocketController) Ws(con *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = con.ReadMessage(); err != nil {
			if err != nil {
				break
			}
			c.netService.OnDisconnect(con)
			break
		}
		c.netService.OnIncomingMessage(con, mt, msg)

	}
}

func (w *WebsocketController) CheckPlayerName(c *fiber.Ctx) error {
	var payload CheckNamePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	game := w.netService.GetGameByCode(payload.Code)
	if game == nil {
		return c.SendStatus(fiber.StatusOK)
	}

	if w.netService.IsNameTakenInGame(payload.Code, payload.Name) {
		return c.SendStatus(fiber.StatusConflict)
	}

	for _, player := range game.Players {
		if player.Name == payload.Name {
			return c.SendStatus(fiber.StatusConflict)
		}
	}
	return c.SendStatus(fiber.StatusOK)
}

func (w *WebsocketController) CheckGamePin(c *fiber.Ctx) error {
	var payload struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	game := w.netService.GetGameByCode(payload.Code)
	if game == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if game.State != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Game has already started"})
	}

	return c.SendStatus(fiber.StatusOK)
}
