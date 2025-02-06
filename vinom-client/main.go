package main

import (
	"crypto/rsa"
	"net"
	"time"

	"github.com/beka-birhanu/vinom-client/config"
	"github.com/beka-birhanu/vinom-client/controller"
	"github.com/beka-birhanu/vinom-client/dmn"
	"github.com/beka-birhanu/vinom-client/infrastruture/crypto"
	"github.com/beka-birhanu/vinom-client/infrastruture/http"
	gamepb "github.com/beka-birhanu/vinom-client/infrastruture/pb_encoder/game"
	udppb "github.com/beka-birhanu/vinom-client/infrastruture/pb_encoder/udp"
	"github.com/beka-birhanu/vinom-client/infrastruture/udp"
	"github.com/beka-birhanu/vinom-client/service"
	"github.com/rivo/tview"
)

var player *dmn.Player
var app *tview.Application

func main() {
	httpClient := http.NewHttpClient(config.Envs.ServerAddr)
	authService, err := service.NewAuth(httpClient, config.Envs.LoginUri, config.Envs.RegisterUri)
	if err != nil {
		panic(err)
	}
	matchService, err := service.NewMatchMaking(service.MatchMakingConfig{
		HttpClient: httpClient,
		MatchUri:   config.Envs.MatchUri,
	})

	app = tview.NewApplication()
	matchPage, err := controller.NewMatchingRoomPage(matchService, startGame)
	if err != nil {
		panic(err)
	}

	authPage, err := controller.NewAuthPage(authService, func(p *dmn.Player, token string) {
		player = p
		err := matchPage.Start(app, player.ID, token)
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}

	err = authPage.Start(app)
	if err != nil {
		panic(err)
	}
}

func startGame(key []byte, addr string) {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	aesKey := []byte{113, 110, 25, 53, 11, 53, 68, 33, 17, 36, 22, 7, 125, 11, 35, 16, 83, 61, 59, 49, 31, 22, 69, 17, 24, 125, 11, 35, 16, 83, 61, 59}
	updClient, err := udp.NewClientServerManager(
		udp.ClientConfig{
			ServerAddr:         serverAddr,
			Encoder:            &udppb.Protobuf{},
			AsymmCrypto:        crypto.NewRSA(&rsa.PrivateKey{}),
			ServerAsymmPubKey:  key,
			SymmCrypto:         crypto.NewAESCBC(),
			ClientSymmKey:      aesKey,
			OnConnectionSucces: func() {},
		},
		udp.ClientWithPingInterval(2*time.Second),
	)

	if err != nil {
		panic(err)
	}

	gameService, err := service.NewGameServer(&service.GameServerConfig{
		ServerConnection: updClient,
		Encoder:          &gamepb.Protobuf{},
		PlayerID:         player.ID,
	})

	if err != nil {
		panic(err)
	}

	gamePage, err := controller.NewGame(gameService, player.ID)
	if err != nil {
		panic(err)
	}

	gamePage.Start(app, player.ID[:])
}
