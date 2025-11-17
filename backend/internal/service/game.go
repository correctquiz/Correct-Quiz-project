package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"CorrectQuiz.com/quiz/internal/entity"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Player struct {
	Id                  uuid.UUID       `json:"id"`
	Name                string          `json:"name"`
	Connection          *websocket.Conn `json:"-"`
	Answered            bool            `json:"-"`
	LastAwardedPoints   int             `json:"-"`
	Points              int             `json:"score"`
	CorrectStreak       int             `json:"correctStreak"`
	AnswerTimeRemaining int             `json:"-"`
	MaxCorrectStreak    int             `json:"-"`
	CurrentAnswer       int             `json:"-"`
}

type GameState int

const (
	LobbyState GameState = iota
	PlayState
	IntermissionState
	RevealState
	EndState
	GameEndedState
)

type LeaderboardEntry struct {
	Name         string `json:"name"`
	Points       int    `json:"points"`
	CorrectCount int    `json:"correctCount"`
}

type Game struct {
	Id                  uuid.UUID
	Quiz                entity.Quiz
	Code                string
	State               GameState
	Time                int
	Players             []*Player
	playersMutex        sync.RWMutex
	CurrentQuestion     int
	Host                *websocket.Conn
	netService          *NetService
	correctAnswerCounts map[uuid.UUID]int
	ctx                 context.Context
	cancelFunc          context.CancelFunc
}

type PlayerAnswerFeedbackPacket struct {
	IsCorrect          bool  `json:"isCorrect"`
	CorrectAnswerIndex []int `json:"correctAnswerIndex"`
	StreakBonus        int   `json:"streakBonus"`
	MaxStreak          int   `json:"maxStreak"`
}

type AnswerReceivedPacket struct {
	PlayerId    uuid.UUID `json:"player_id"`
	ChoiceIndex int       `json:"choice_index"`
}

func generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

func (g *Game) AddPlayer(player *Player) {
	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()
	g.Players = append(g.Players, player)
}

func (g *Game) RemovePlayer(playerId uuid.UUID) {
	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	var removedPlayerName string
	var newPlayers []*Player
	for _, p := range g.Players {
		if p.Id != playerId {
			newPlayers = append(newPlayers, p)
		} else {
			removedPlayerName = p.Name
		}
	}
	if len(newPlayers) < len(g.Players) {
		g.Players = newPlayers
		log.Printf("Game %s: Player %s (ID: %s) removed. Total players: %d", g.Code, removedPlayerName, playerId, len(g.Players))
	}
}

func newGame(quiz entity.Quiz, host *websocket.Conn, netService *NetService) Game {
	ctx, cancel := context.WithCancel(context.Background())
	return Game{
		Id:                  uuid.New(),
		Quiz:                quiz,
		Code:                generateCode(),
		Players:             []*Player{},
		State:               LobbyState,
		Time:                60,
		Host:                host,
		netService:          netService,
		correctAnswerCounts: make(map[uuid.UUID]int),
		ctx:                 ctx,
		cancelFunc:          cancel,
	}
}

func (g *Game) Start() {
	g.ChangeState(PlayState)
	g.CurrentQuestion = 0

	if len(g.Quiz.Questions) == 0 {
		fmt.Println("Error: Quiz has no questions, cannot start game.")
		return
	}

	g.Time = g.Quiz.Questions[0].Time

	questionPacket := QuestionShowPacket{
		Question:      g.Quiz.Questions[0],
		QuestionIndex: g.CurrentQuestion,
	}

	g.BroadcastPacket(questionPacket, true)

	go func(gameCtx context.Context) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				g.Tick()
			case <-gameCtx.Done():
				return
			}
		}
	}(g.ctx)
}

func (g *Game) ResetPlayerAnswerStates() {
	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()
	for _, player := range g.Players {
		player.Answered = false
	}
}

func (g *Game) NextQuestion() {
	g.CurrentQuestion++

	g.ResetPlayerAnswerStates()
	g.ChangeState(PlayState)

	g.Time = g.Quiz.Questions[g.CurrentQuestion].Time

	g.BroadcastPacket(QuestionShowPacket{
		Question:      g.Quiz.Questions[g.CurrentQuestion],
		QuestionIndex: g.CurrentQuestion,
	}, true)
}

func (g *Game) Reveal() {
	g.Time = 10
	g.playersMutex.RLock()
	type PlayerRevealInfo struct {
		Connection *websocket.Conn
		Points     int
	}
	reveals := make([]PlayerRevealInfo, 0, len(g.Players))
	for _, player := range g.Players {
		if player.Connection != nil {
			reveals = append(reveals, PlayerRevealInfo{
				Connection: player.Connection,
				Points:     player.LastAwardedPoints,
			})
		}
	}

	g.playersMutex.RUnlock()

	for _, player := range g.Players {
		g.netService.SendPacket(player.Connection, PlayerRevealPacket{
			Points: player.LastAwardedPoints,
		})
	}
	g.ChangeState(RevealState)
}

func (g *Game) Tick() {
	if g.Time > 0 {
		g.Time--
		g.netService.SendPacket(g.Host, TickPacket{
			Tick: g.Time,
		})
	}

	if g.Time == 0 {
		switch g.State {
		case PlayState:
			{
				g.State = RevealState

				g.netService.SendPacket(g.Host, ChangeGameStatePacket{
					State: RevealState,
				})

				g.sendHostReveal(g.CurrentQuestion)
				g.sendPlayerResults()
			}
		}
	}
}

func (g *Game) GetCorrectAnswerCount(playerID uuid.UUID) (int, bool) {
	if g.correctAnswerCounts == nil {
		return 0, false
	}
	count, ok := g.correctAnswerCounts[playerID]
	return count, ok
}

func (g *Game) Intermission() {
	if g.CurrentQuestion >= len(g.Quiz.Questions)-1 {
		g.ChangeState(EndState)

		leaderboardData := g.getLeaderboard()
		leaderboardPacket := LeaderboardPacket{
			Points: leaderboardData,
		}

		g.playersMutex.RLock()
		type PlayerRankInfo struct {
			Connection *websocket.Conn
			Rank       int
			Name       string
		}
		ranks := make([]PlayerRankInfo, 0, len(g.Players))

		for i, entry := range leaderboardData {
			for _, player := range g.Players {
				if player.Name == entry.Name {
					rankPacket := PlayerRankPacket{Rank: i + 1}
					g.netService.SendPacket(player.Connection, rankPacket)
					break
				}
			}
		}
		hostConnection := g.Host
		g.playersMutex.RUnlock()

		g.netService.SendPacket(g.Host, leaderboardPacket)
		for _, rankInfo := range ranks {
			rankPacket := PlayerRankPacket{Rank: rankInfo.Rank}
			g.netService.SendPacket(rankInfo.Connection, rankPacket)
		}
		g.netService.SendPacket(hostConnection, leaderboardPacket)

	} else {
		g.State = IntermissionState

		hostStatePacket := ChangeGameStatePacket{State: IntermissionState}
		leaderboardData := g.getLeaderboard()
		leaderboardPacket := LeaderboardPacket{Points: leaderboardData}

		g.netService.SendPacket(g.Host, hostStatePacket)
		g.netService.SendPacket(g.Host, leaderboardPacket)
	}
}

func (g *Game) getLeaderboard() []LeaderboardEntry {

	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	sort.Slice(g.Players, func(i, j int) bool {
		if g.Players[i].Points != g.Players[j].Points {
			return g.Players[i].Points > g.Players[j].Points
		}
		iCount := g.correctAnswerCounts[g.Players[i].Id]
		jCount := g.correctAnswerCounts[g.Players[j].Id]
		return iCount > jCount
	})

	leaderboard := []LeaderboardEntry{}
	for _, player := range g.Players {
		leaderboard = append(leaderboard, LeaderboardEntry{
			Name:         player.Name,
			Points:       player.Points,
			CorrectCount: g.correctAnswerCounts[player.Id],
		})
	}

	return leaderboard
}

func (g *Game) ChangeState(state GameState) {
	g.State = state
	g.BroadcastPacket(ChangeGameStatePacket{
		State: state,
	}, true)
}

func (g *Game) BroadcastPacket(packet any, includeHost bool) error {
	g.playersMutex.RLock()

	connections := make([]*websocket.Conn, 0, len(g.Players))
	for _, player := range g.Players {
		if player.Connection != nil {
			connections = append(connections, player.Connection)
		}
	}

	g.playersMutex.RUnlock()

	for _, conn := range connections {
		g.netService.SendPacket(conn, packet)
	}

	if includeHost && g.Host != nil {
		err := g.netService.SendPacket(g.Host, packet)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) OnPlayerJoin(name string, connection *websocket.Conn) {
	player := Player{
		Id:            uuid.New(),
		Name:          name,
		Connection:    connection,
		Points:        0,
		CorrectStreak: 0,
	}
	g.AddPlayer(&player)

	g.playersMutex.RLock()
	hostConnection := g.Host
	g.playersMutex.RUnlock()

	g.netService.SendPacket(connection, ChangeGameStatePacket{
		State: g.State,
	})
	g.netService.SendPacket(hostConnection, PlayerJoinPacket{
		Player: player,
	})
	g.netService.SendPacket(connection, PlayerJoinPacket{
		Player: player,
	})
}

func (g *Game) KickPlayer(playerID string) error {
	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	var playerToKick *Player = nil
	var playerIndex = -1

	for i, player := range g.Players {
		if player.Id.String() == playerID {
			playerToKick = player
			playerIndex = i
			break
		}
	}

	if playerToKick == nil {
		return errors.New("player not found")
	}

	if playerToKick.Connection != nil {
		playerToKick.Connection.WriteMessage(websocket.CloseMessage, []byte{})
	}

	g.Players = append(g.Players[:playerIndex], g.Players[playerIndex+1:]...)

	return nil
}

func (g *Game) sendHostReveal(questionIndex int) {
	if questionIndex < 0 || questionIndex >= len(g.Quiz.Questions) {
		return
	}

	currentQuestion := g.Quiz.Questions[questionIndex]
	var correctAnswerIndex []int
	for i, choice := range currentQuestion.Choices {
		if choice.Correct {
			correctAnswerIndex = append(correctAnswerIndex, i)
		}
	}

	counts := make([]int, len(currentQuestion.Choices))

	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()

	for _, player := range g.Players {
		if player.Answered {
			playerChoiceIndex := player.CurrentAnswer
			if playerChoiceIndex >= 0 && playerChoiceIndex < len(counts) {
				counts[playerChoiceIndex]++
			}
		}
	}

	packet := QuestionRevealPacket{
		Question:           currentQuestion,
		CorrectAnswerIndex: correctAnswerIndex,
		AnswerCounts:       counts,
	}
	g.netService.SendPacket(g.Host, packet)
}

func (g *Game) OnPlayerAnswer(questionIndex int, choiceIndex int, player *Player) {

	g.playersMutex.Lock()

	if player.Answered {
		g.playersMutex.Lock()
		return
	}

	player.Answered = true
	player.CurrentAnswer = choiceIndex
	player.AnswerTimeRemaining = g.Time

	allAnswered := true
	if len(g.Players) == 0 {
		allAnswered = false
	}

	for _, p := range g.Players {
		if !p.Answered {
			allAnswered = false
			break
		}
	}

	if allAnswered && g.State == PlayState {
		g.Time = 0
	}

	g.playersMutex.Unlock()

	if questionIndex < 0 || questionIndex >= len(g.Quiz.Questions) {
		return
	}

	currentQuestion := g.Quiz.Questions[questionIndex]
	if choiceIndex < 0 || choiceIndex >= len(currentQuestion.Choices) {
		return
	}
}

func (g *Game) sendPlayerResults() {
	if g.CurrentQuestion < 0 || g.CurrentQuestion >= len(g.Quiz.Questions) {
		return
	}

	currentQuestion := g.Quiz.Questions[g.CurrentQuestion]
	var correctAnswerIndex []int
	for i, choice := range currentQuestion.Choices {
		if choice.Correct {
			correctAnswerIndex = append(correctAnswerIndex, i)
		}
	}

	totalTime := float64(currentQuestion.Time)
	if totalTime <= 0 {
		totalTime = 60.0
	}

	const basePointsForCorrect = 100.0

	type PlayerPacketPair struct {
		Connection *websocket.Conn
		Feedback   PlayerAnswerFeedbackPacket
		Reveal     PlayerRevealPacket
	}
	packetsToSend := []PlayerPacketPair{}
	correctPlayers := []*Player{}

	for _, player := range g.Players {
		if player.Answered {
			isCorrect := false
			for _, correctIdx := range correctAnswerIndex {
				if player.CurrentAnswer == correctIdx {
					isCorrect = true
					break
				}
			}
			if isCorrect {
				correctPlayers = append(correctPlayers, player)
			}
		}
	}

	sort.Slice(correctPlayers, func(i, j int) bool {
		return correctPlayers[i].AnswerTimeRemaining > correctPlayers[j].AnswerTimeRemaining
	})

	pointsMap := make(map[uuid.UUID]int)
	rankedScores := []int{26, 24, 22}
	baseCorrectScore := 20

	for i, player := range correctPlayers {
		score := baseCorrectScore
		if i < len(rankedScores) {
			score = rankedScores[i]
		}
		pointsMap[player.Id] = score
	}

	g.playersMutex.Lock()
	for _, player := range g.Players {
		isCorrect := false
		if player.Answered {
			for _, correctIdx := range correctAnswerIndex {
				if player.CurrentAnswer == correctIdx {
					isCorrect = true
					break
				}
			}
		}

		var awardedPointsThisRound int = 0
		var streakBonus int = 0

		if isCorrect {
			awardedPointsThisRound = int(basePointsForCorrect)
			awardedPointsThisRound += pointsMap[player.Id]
			g.correctAnswerCounts[player.Id]++

			player.CorrectStreak++

			if player.CorrectStreak > player.MaxCorrectStreak {
				player.MaxCorrectStreak = player.CorrectStreak
			}

			if player.CorrectStreak == 2 {
				streakBonus = 10
				awardedPointsThisRound += streakBonus
			} else if player.CorrectStreak == 3 {
				streakBonus = 20
				awardedPointsThisRound += streakBonus
				player.CorrectStreak = 0
			}

			player.Points += awardedPointsThisRound

		} else {
			player.CorrectStreak = 0
		}

		player.LastAwardedPoints = awardedPointsThisRound

		feedbackPacket := PlayerAnswerFeedbackPacket{
			IsCorrect:          isCorrect,
			CorrectAnswerIndex: correctAnswerIndex,
			StreakBonus:        streakBonus,
			MaxStreak:          player.MaxCorrectStreak,
		}
		revealPacket := PlayerRevealPacket{
			Points: player.Points,
		}

		if player.Connection != nil {
			packetsToSend = append(packetsToSend, PlayerPacketPair{
				Connection: player.Connection,
				Feedback:   feedbackPacket,
				Reveal:     revealPacket,
			})
		}

	}
	g.playersMutex.Unlock()

	for _, pair := range packetsToSend {
		if pair.Connection != nil {
			g.netService.SendPacket(pair.Connection, pair.Feedback)
			g.netService.SendPacket(pair.Connection, pair.Reveal)
		}
	}
}
