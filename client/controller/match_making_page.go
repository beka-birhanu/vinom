package controller

import (
	"strings"
	"time"

	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/gdamore/tcell/v2"
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
	footer := tview.NewTextView().
		SetText("").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	footer.SetBackgroundColor(catBase)

	form := tview.NewForm()
	form.SetBackgroundColor(catBase)
	form.SetButtonBackgroundColor(tcell.GetColor("#a6e3a1")) // Catppuccin Green
	form.SetButtonTextColor(catBase)
	form.AddButton("Find Match", func() {
		done := make(chan struct{})

		// Animate footer with dots while searching
		go func() {
			dots := 0
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					dots = (dots % 3) + 1
					dotStr := strings.Repeat(".", dots)
					app.QueueUpdateDraw(func() {
						footer.SetText("[#f9e2af]Searching" + dotStr + "[#cdd6f4]")
					})
				case <-done:
					return
				}
			}
		}()

		go func() {
			pubKey, addr, err := m.matchService.Match(ID, token)
			close(done)
			if err != nil {
				app.QueueUpdateDraw(func() {
					footer.SetText("[#f38ba8]" + err.Error() + "[#cdd6f4]")
				})
				return
			}
			app.QueueUpdateDraw(func() {
				footer.SetText("[#a6e3a1]Found a match![#cdd6f4]")
			})
			m.onMatch(pubKey, addr)
		}()
	})

	form.AddButton("Cancel", func() {
		app.Stop()
	})

	frame := tview.NewFrame(form).SetBorders(1, 1, 0, 0, 1, 1)
	frame.SetBorder(true)
	frame.SetTitle(" Matchmaking ")
	frame.SetBorderColor(catMauve)
	frame.SetTitleColor(catBlue)
	frame.SetBackgroundColor(catBase)

	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, true).
		AddItem(footer, 1, 0, false)
	content.SetBackgroundColor(catBase)

	centered := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(content, 0, 2, true).
			AddItem(nil, 0, 1, false), 0, 2, true).
		AddItem(nil, 0, 1, false)
	centered.SetBackgroundColor(catBase)

	return centered
}
