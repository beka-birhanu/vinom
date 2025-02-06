package i

// ClientManager defines the interface for the UDP client socket manager.
type ClientManager interface {
	// Connect establishes the initial handshake with the server.
	Connect([]byte) error

	// Disconnect closes the connection and cleans up resources.
	Disconnect()

	// SendToServer encrypts and sends a message of the specified type to the server.
	SendToServer(t byte, message []byte) error

	// SetOnServerResponse updates onServerResponse func.
	SetOnServerResponse(f func(byte, []byte))

	// SetOnPingResult updates onServerResponse func.
	SetOnPingResult(f func(int64))
}
