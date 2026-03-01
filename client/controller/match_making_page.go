package controller

import (
	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

type matchHandler func([]byte, string)

type MatchingRoomPage struct {
	matchService i.MatchMaker
	onMatch      matchHandler
}

func NewMatchingRoomPage(ms i.MatchMaker, onMatch matchHandler) (*MatchingRoomPage, error) {
	return &MatchingRoomPage{
		matchService: ms,
		onMatch:      onMatch,
	}, nil
}

func (m *MatchingRoomPage) Start(app *tview.Application, ID uuid.UUID, token string) error {
	if err := app.SetRoot(m.matchingRoomUI(app, ID, token), true).Run(); err != nil {
		return err
	}
	return nil
}

func (m *MatchingRoomPage) matchingRoomUI(app *tview.Application, ID uuid.UUID, token string) tview.Primitive {
	header := tview.NewTextView().SetText("Matching Room").SetTextAlign(tview.AlignCenter)
	footer := tview.NewTextView().SetText("").SetTextAlign(tview.AlignLeft)

	form := tview.NewForm()
	form.AddButton("Find Match", func() {
		footer.SetText("Searching for match...")
		go func(footer *tview.TextView, ID uuid.UUID) {
			pubKey, addr, err := m.matchService.Match(ID, token)
			if err != nil {
				footer.SetText(err.Error())
				app.Draw()
				return
			}

			m.onMatch(pubKey, addr)
			footer.SetText("Found a match for you!")
			app.Draw()
		}(footer, ID)
	})

	form.AddButton("Cancel", func() {
		app.Stop()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(footer, 0, 1, false)

	return flex
}
