package service

import (
	"sync"

	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/google/uuid"
)

const (
	moveActionType         = 3 << iota // Action type for movement.
	stateRequestActionType             // Action type for state requests.

	gameStateRecordType = 10
	gameEndedRecordType = 11
)

type GameServer struct {
	serverConnection i.ClientManager
	encoder          i.GameEncoder
	onGameEnd        func(i.GameState)
	gameState        i.GameState
	playerID         uuid.UUID
	onStateChange    func(i.GameState)
	onPingResult     func(int64)
	sync.Mutex
}

type GameServerConfig struct {
	ServerConnection i.ClientManager
	Encoder          i.GameEncoder
	OnGameEnd        func(i.GameState)
	PlayerID         uuid.UUID
}

func NewGameServer(cfg *GameServerConfig) (i.GameServer, error) {
	server := &GameServer{
		serverConnection: cfg.ServerConnection,
		encoder:          cfg.Encoder,
		playerID:         cfg.PlayerID,
	}

	server.serverConnection.SetOnServerResponse(server.handleServerResponse)
	server.serverConnection.SetOnPingResult(server.handlePingResponse)
	return server, nil
}

func (g *GameServer) Start(authToken []byte) error {
	err := g.serverConnection.Connect(authToken)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameServer) Stop() error {
	g.serverConnection.Disconnect()
	return nil
}

// move implements i.GameServer.
func (g *GameServer) Move(direction string) {
	action := g.encoder.NewAction()
	action.SetDirection(direction)
	action.SetID(g.playerID)
	action.SetFrom(g.playerPosition())

	payload, err := g.encoder.MarshalAction(action)
	if err != nil {
		return
	}

	err = g.serverConnection.SendToServer(moveActionType, payload)
	if err != nil {
		return
	}
}

func (g *GameServer) handleServerResponse(t byte, p []byte) {
	g.Lock()
	defer g.Unlock()

	gameState, err := g.encoder.UnmarshalGameState(p)
	if err != nil {
		return
	}

	if t == gameEndedRecordType {
		g.onGameEnd(gameState)
		return
	}

	if g.gameState == nil || g.gameState.GetVersion() < gameState.GetVersion() {
		g.gameState = gameState
		g.onStateChange(g.gameState)
	}
}

func (g *GameServer) handlePingResponse(ping int64) {
	g.onPingResult(ping)
}

func (g *GameServer) playerPosition() i.CellPosition {
	g.Lock()
	defer g.Unlock()

	for _, player := range g.gameState.RetrivePlayers() {
		if player.GetID() == g.playerID {
			return player.RetrivePos()
		}
	}
	return nil // code will no reach this; or at least I hope it does not
}

func (g *GameServer) SetOnStateChange(f func(i.GameState)) {
	g.onStateChange = f
}

func (g *GameServer) SetOnPingResult(f func(int64)) {
	g.onPingResult = f
}
