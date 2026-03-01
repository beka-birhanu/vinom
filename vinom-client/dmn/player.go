package dmn

import "github.com/google/uuid"

type Player struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Rating   int       `json:"rating"`
}
