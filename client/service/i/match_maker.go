package i

import (
	"github.com/google/uuid"
)

type MatchMaker interface {
	Match(ID uuid.UUID, token string) ([]byte, string, error)
}
