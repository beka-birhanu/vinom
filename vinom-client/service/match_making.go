package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/google/uuid"
)

type MatchMaking struct {
	httpClient i.HttpRequester
	matchUri   string
}

type MatchMakingConfig struct {
	HttpClient i.HttpRequester
	MatchUri   string
}

func NewMatchMaking(mc MatchMakingConfig) (*MatchMaking, error) {
	return &MatchMaking{
		httpClient: mc.HttpClient,
		matchUri:   mc.MatchUri,
	}, nil
}

func (mm *MatchMaking) Match(ID uuid.UUID, token string) ([]byte, string, error) {
	sentAt := time.Now().UnixNano() / int64(time.Millisecond)
	body := MatchRequest{ID: ID, SentAt: sentAt}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, "", err
	}

	_, err = mm.httpClient.Post(mm.matchUri, bytes.NewReader(payload), token)
	if err != nil {
		return nil, "", err
	}

	maxTime := time.NewTimer(time.Minute)
	matchInfoUri := fmt.Sprintf("%s/%s", mm.matchUri, ID)
	for {
		select {
		case <-maxTime.C:
			return nil, "", errors.New("Match Request Timeout.")
		default:
			response, err := mm.httpClient.Get(matchInfoUri, token)
			if err != nil {
				time.Sleep(2 * time.Second)
			} else {
				return parseInfoResponse(response)
			}
		}
	}
}

func parseInfoResponse(response io.Reader) ([]byte, string, error) {
	payload, err := io.ReadAll(response)
	if err != nil {
		return nil, "", err
	}

	var matchInfo MatchInfoResponse
	err = json.Unmarshal(payload, &matchInfo)
	if err != nil {
		return nil, "", err
	}

	return matchInfo.SocketPubKey, matchInfo.SocketAddr, nil
}
