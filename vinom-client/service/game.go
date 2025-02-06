package service

import (
	"sync"

	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/google/uuid"
)

const (
	moveActionType         = 1 << iota // Action type for movement.
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
	stateChan        chan i.GameState
	pingChan         chan int64
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

	g.stateChan = make(chan i.GameState)
	g.pingChan = make(chan int64)
	return g.requestGameState()
}

func (g *GameServer) Stop() error {
	g.serverConnection.Disconnect()
	close(g.stateChan)
	close(g.pingChan)
	return nil
}

func (g *GameServer) requestGameState() error {
	return g.serverConnection.SendToServer(stateRequestActionType, []byte{})
}

// move implements i.GameServer.
func (g *GameServer) Move(direction string) {
	action := g.encoder.NewAction()
	action.SetDirection(direction)
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

// pingChan implements i.GameServer.
func (g *GameServer) PingChan() <-chan int64 {
	return g.pingChan
}

// stateChan implements i.GameServer.
func (g *GameServer) StateChan() <-chan i.GameState {
	return g.stateChan
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

	if g.gameState.GetVersion() < gameState.GetVersion() {
		g.stateChan <- gameState
		g.gameState = gameState
	}
}

func (g *GameServer) handlePingResponse(ping int64) {
	g.pingChan <- ping
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
