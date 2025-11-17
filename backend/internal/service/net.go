package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"CorrectQuiz.com/quiz/internal/entity"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type NetService struct {
	quizService *QuizService
	games       []*Game
	gamesMutex  sync.RWMutex
}

func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		games:       []*Game{},
		gamesMutex:  sync.RWMutex{},
	}
}

type ConnectPacket struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type HostGamePacket struct {
	QuizId string `json:"quizId"`
}

type QuestionShowPacket struct {
	Question      entity.QuizQuestion `json:"question"`
	QuestionIndex int                 `json:"questionIndex"`
}

type ChangeGameStatePacket struct {
	State GameState `json:"state"`
	Code  string    `json:"code,omitempty"`
}

type PlayerJoinPacket struct {
	Player Player `json:"player"`
}

type StartGamePacket struct {
}

type TickPacket struct {
	Tick int `json:"tick"`
}

type QuestionAnswerPacket struct {
	Question int `json:"question"`
	Choice   int `json:"choice"`
}

type PlayerRevealPacket struct {
	Points int `json:"points"`
}

type QuestionRevealPacket struct {
	Question           entity.QuizQuestion `json:"question"`
	CorrectAnswerIndex []int               `json:"correctAnswerIndex"`
	AnswerCounts       []int               `json:"answerCounts"`
}

type LeaderboardPacket struct {
	Points []LeaderboardEntry `json:"points"`
}

type PlayerLeavePacket struct {
	PlayerId uuid.UUID `json:"playerId"`
}

type PlayerRankPacket struct {
	Rank int `json:"rank"`
}

type KickPlayerPacket struct {
	PlayerId string `json:"playerId"`
}

type HostLeavePacket struct{}

type NextQuestionPacket struct{}

func (c *NetService) packetIdtoPacket(packetId uint8) any {
	switch packetId {
	case 0:
		{
			return &ConnectPacket{}
		}
	case 1:
		{
			return &HostGamePacket{}
		}
	case 3:
		{
			return &ChangeGameStatePacket{}
		}
	case 5:
		{
			return &StartGamePacket{}
		}
	case 7:
		{
			return &QuestionAnswerPacket{}
		}
	case 8:
		{
			return &PlayerAnswerFeedbackPacket{
				IsCorrect:          false,
				CorrectAnswerIndex: nil,
			}
		}
	case 9:
		{
			return &AnswerReceivedPacket{}
		}
	case 13:
		{
			return &NextQuestionPacket{}
		}
	case 14:
		{
			return &PlayerRankPacket{}
		}
	case 15:
		{
			return &KickPlayerPacket{}
		}
	case 16:
		{
			return &HostLeavePacket{}
		}
	case 17:
		{
			return &PlayerLeavePacket{}
		}
	}

	return nil
}

func (c *NetService) packettoPacketId(packet any) (uint8, error) {
	switch packet.(type) {
	case HostGamePacket:
		{
			return 1, nil
		}
	case QuestionShowPacket:
		{
			return 2, nil
		}
	case ChangeGameStatePacket:
		{
			return 3, nil
		}
	case PlayerJoinPacket:
		{
			return 4, nil
		}
	case TickPacket:
		{
			return 6, nil
		}
	case PlayerAnswerFeedbackPacket:
		{
			return 8, nil
		}
	case QuestionRevealPacket:
		{
			return 10, nil
		}
	case PlayerRevealPacket:
		{
			return 11, nil
		}
	case LeaderboardPacket:
		{
			return 12, nil
		}
	case PlayerRankPacket:
		{
			return 14, nil
		}
	case *HostLeavePacket:
		{
			return 16, nil
		}
	case PlayerLeavePacket:
		{
			return 17, nil
		}
	}
	return 0, errors.New("invalid packet type")
}

func (c *NetService) GetGameByCode(code string) *Game {
	c.gamesMutex.RLock()
	defer c.gamesMutex.RUnlock()
	for _, game := range c.games {
		if game.Code == code {
			return game
		}
	}
	return nil
}

func (c *NetService) GetGameByHost(host *websocket.Conn) *Game {
	c.gamesMutex.RLock()
	defer c.gamesMutex.RUnlock()
	for _, game := range c.games {
		if game.Host == host {
			return game
		}
	}
	return nil
}

func (c *NetService) GetGameByPlayer(con *websocket.Conn) (*Game, *Player) {
	c.gamesMutex.RLock()
	defer c.gamesMutex.RUnlock()
	for _, game := range c.games {
		game.playersMutex.RLock()
		for _, player := range game.Players {
			if player.Connection == con {
				game.playersMutex.RUnlock()
				return game, player
			}
		}
		game.playersMutex.RUnlock()
	}
	return nil, nil
}

func (c *NetService) IsNameTakenInGame(code string, name string) bool {

	game := c.GetGameByCode(code)
	if game == nil {
		return false
	}

	game.playersMutex.RLock()
	defer game.playersMutex.RUnlock()

	for _, player := range game.Players {
		if player.Name == name {
			return true
		}
	}

	return false
}

func (c *NetService) FindActiveGameByCode(code string) (*Game, error) {
	c.gamesMutex.RLock()
	defer c.gamesMutex.RUnlock()

	for _, game := range c.games {
		if game.Code == code {
			return game, nil
		}
	}
	return nil, errors.New("active game not found")
}

func (c *NetService) handleHostLeave(con *websocket.Conn) {
	game := c.GetGameByHost(con)

	if game != nil {
		endPacket := ChangeGameStatePacket{
			State: GameEndedState,
		}

		if game.cancelFunc != nil {
			game.cancelFunc()
		}

		for _, p := range game.Players {
			if p.Connection != nil {
				err := c.SendPacket(p.Connection, endPacket)
				if err != nil {
				}

				p.Connection.Close()
			}
		}
		c.gamesMutex.Lock()
		var newGames []*Game
		for _, g := range c.games {
			if g.Id != game.Id {
				newGames = append(newGames, g)
			}
		}
		c.games = newGames
		c.gamesMutex.Unlock()
	}
}

func (c *NetService) handlePlayerLeave(con *websocket.Conn, playerId uuid.UUID) {
	game, player := c.GetGameByPlayerId(playerId)
	if game != nil && player != nil {

		game.RemovePlayer(player.Id)

		leavePacket := PlayerLeavePacket{
			PlayerId: player.Id,
		}
		err := c.SendPacket(game.Host, leavePacket)
		if err != nil {
		}

	}
}

func (c *NetService) OnIncomingMessage(con *websocket.Conn, mt int, msg []byte) {

	if len(msg) < 2 {
		return
	}

	packetId := msg[0]
	data := msg[1:]

	packet := c.packetIdtoPacket(packetId)
	if packet == nil {
		return
	}

	err := json.Unmarshal(data, packet)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch data := packet.(type) {
	case *ConnectPacket:
		{
			game := c.GetGameByCode(data.Code)
			if game == nil {
				return
			}
			game.OnPlayerJoin(data.Name, con)
			break
		}

	case *HostGamePacket:
		{

			id64, err := strconv.ParseUint(data.QuizId, 10, 64)
			if err != nil {
				return
			}
			quizId := uint(id64)

			quiz, err := c.quizService.GetQuizById(quizId)
			if err != nil {
				fmt.Println(err)
				return
			}

			if quiz == nil {
				fmt.Println("Quiz not found (ID:", quizId, ")")
				return
			}

			game := newGame(*quiz, con, c)
			c.gamesMutex.Lock()
			c.games = append(c.games, &game)
			c.gamesMutex.Unlock()

			c.SendPacket(con, HostGamePacket{
				QuizId: game.Code,
			})

			c.SendPacket(con, ChangeGameStatePacket{
				State: game.State,
				Code:  game.Code,
			})
			break
		}

	case *StartGamePacket:
		{
			game := c.GetGameByHost(con)
			if game == nil {
				return
			}
			game.Start()
			break
		}
	case *KickPlayerPacket:
		{
			game := c.GetGameByHost(con)
			if game == nil {
				return
			}

			var playerToKick *Player = nil
			for _, player := range game.Players {
				if player.Id.String() == data.PlayerId {
					playerToKick = player
					break
				}
			}

			if playerToKick == nil {
				return
			}

			game.KickPlayer(data.PlayerId)

			break
		}
	case *QuestionAnswerPacket:
		{
			game, player := c.GetGameByPlayer(con)
			if game == nil {
				return
			}

			game.OnPlayerAnswer(data.Question, data.Choice, player)
			break

		}

	case *ChangeGameStatePacket:
		{
			game := c.GetGameByHost(con)
			if game == nil {
				return
			}

			if data.State == IntermissionState {
				game.Intermission()
			}
			break
		}
	case *NextQuestionPacket:
		{
			game := c.GetGameByHost(con)
			if game == nil {
				return
			}
			game.NextQuestion()
			break
		}
	case *PlayerLeavePacket:
		{
			c.handlePlayerLeave(con, data.PlayerId)
			break
		}
	case *HostLeavePacket:
		{
			c.handleHostLeave(con)
			break
		}
	}
}

func (c *NetService) GetGameByPlayerId(playerId uuid.UUID) (*Game, *Player) {
	c.gamesMutex.RLock()
	defer c.gamesMutex.RUnlock()

	for _, game := range c.games {
		game.playersMutex.RLock()
		for _, player := range game.Players {
			if player.Id == playerId {
				game.playersMutex.RUnlock()
				return game, player
			}
		}
		game.playersMutex.RUnlock()
	}
	return nil, nil
}

func (c *NetService) OnDisconnect(con *websocket.Conn) {
	game, player := c.GetGameByPlayer(con)
	if game != nil && player != nil {

		game.RemovePlayer(player.Id)
		leavePacket := PlayerLeavePacket{PlayerId: player.Id}

		err := c.SendPacket(game.Host, leavePacket)
		if err != nil {
		}
		return
	}

	game = c.GetGameByHost(con)
	if game != nil {

		if game.cancelFunc != nil {
			game.cancelFunc()
		}

		endPacket := ChangeGameStatePacket{
			State: GameEndedState,
		}

		game.playersMutex.RLock()
		for _, p := range game.Players {
			if p.Connection != nil {
				c.SendPacket(p.Connection, endPacket)
				p.Connection.Close()
			}
		}
		game.playersMutex.RUnlock()

		c.gamesMutex.Lock()
		var newGames []*Game

		for _, g := range c.games {
			if g.Code != game.Code {
				newGames = append(newGames, g)
			}
		}
		c.games = newGames
		c.gamesMutex.Unlock()

		return
	}
}

func (c *NetService) SendPacket(connection *websocket.Conn, packet any) error {
	bytes, err := c.PacketToBytes(packet)
	if err != nil {
		return err
	}
	return connection.WriteMessage(websocket.BinaryMessage, bytes)
}

func (c *NetService) PacketToBytes(packet any) ([]byte, error) {
	packetId, err := c.packettoPacketId(packet)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	final := append([]byte{packetId}, bytes...)
	return final, nil
}
