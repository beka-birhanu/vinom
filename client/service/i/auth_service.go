package i

import "github.com/beka-birhanu/vinom-client/dmn"

type AuthServer interface {
	Login(username string, password string) (*dmn.Player, string, error)
	Register(username string, password string) error
}
