package i

type GameServer interface {
	Move(string)
	Start([]byte) error
	SetOnStateChange(f func(GameState))
	SetOnPingResult(f func(int64))
}
