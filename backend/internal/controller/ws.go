package controller

import (
	"log"
	"time"

	"CorrectQuiz.com/quiz/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gorilla "github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
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
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			con.Close()
		}()
		for {
			<-ticker.C
			if err := con.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	var (
		mt  int
		msg []byte
		err error
	)

	con.SetReadDeadline(time.Now().Add(pongWait))

	con.SetPongHandler(func(appData string) error {
		con.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if mt, msg, err = con.ReadMessage(); err != nil {
			if gorilla.IsCloseError(err) {
				log.Printf("🚨 WS DISCONNECT REASON: Code %d", err.(*gorilla.CloseError).Code)
			} else {
				log.Printf("🚨 WS UNHANDLED ERROR: %v", err)
			}
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			}

			c.netService.OnDisconnect(con)
			break
		}
		log.Printf("✅ [Backend] Received Message! Type: %d, Size: %d", mt, len(msg))
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
