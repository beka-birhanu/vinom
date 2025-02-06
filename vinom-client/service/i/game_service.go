package i

type GameServer interface {
	Move(string)
	Start([]byte) error
	StateChan() <-chan GameState
	PingChan() <-chan int64
}
