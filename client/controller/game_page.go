package controller

import (
	"fmt"
	"sort"
	"strings"

	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// directions maps movement directions (North, South, East, West) to row and column deltas.
var directions = map[tcell.Key]string{
	tcell.KeyUp:    "North",
	tcell.KeyDown:  "South",
	tcell.KeyLeft:  "West",
	tcell.KeyRight: "East",
}

// Additional handling for Vim motions
var vimDirections = map[rune]string{
	'k': "North", // Vim motion: k for up
	'j': "South", // Vim motion: j for down
	'h': "West",  // Vim motion: h for left
	'l': "East",  // Vim motion: l for right
}

type playerDisplay struct {
	maze       string // 2 visible chars + color tags, background reset at end
	scoreboard string // readable label with color tags
}

// Game holds the maze, score, ping, and player information
type Game struct {
	gameServer     i.GameServer
	playerDisplays map[uuid.UUID]playerDisplay
	playerID       uuid.UUID
	app            *tview.Application
	mazeTV         *tview.TextView
	scoreTV        *tview.Table
	pingTV         *tview.TextView
	stopChan       chan struct{}
}

// NewGame creates a new MazeGame instance
func NewGame(gmSrvr i.GameServer, pID uuid.UUID) (*Game, error) {
	mazeTV := tview.NewTextView().SetDynamicColors(true)
	mazeTV.SetBackgroundColor(tcell.GetColor("#11111b")) // Catppuccin Crust
	mazeTV.SetBorder(true)
	mazeTV.SetTitle(" VINOM ")
	mazeTV.SetBorderColor(catMauve)
	mazeTV.SetTitleColor(catBlue)

	scoreTV := tview.NewTable()
	scoreTV.SetBackgroundColor(catBase)
	scoreTV.SetBorder(true)
	scoreTV.SetTitle(" Scoreboard ")
	scoreTV.SetBorderColor(catMauve)
	scoreTV.SetTitleColor(catBlue)

	pingTV := tview.NewTextView().SetDynamicColors(true)
	pingTV.SetBackgroundColor(catBase)
	pingTV.SetBorder(true)
	pingTV.SetTitle(" Network ")
	pingTV.SetBorderColor(catMauve)
	pingTV.SetTitleColor(catBlue)

	return &Game{
		gameServer: gmSrvr,
		playerID:   pID,
		mazeTV:     mazeTV,
		scoreTV:    scoreTV,
		pingTV:     pingTV,
		stopChan:   make(chan struct{}),
	}, nil
}

func (g *Game) handleInput(event *tcell.EventKey) *tcell.EventKey {
	if direction, ok := directions[event.Key()]; ok {
		g.gameServer.Move(direction)
	} else if direction, ok := vimDirections[event.Rune()]; ok {
		g.gameServer.Move(direction)
	} else if event.Key() == tcell.KeyCtrlC {
		g.stopChan <- struct{}{}
	}
	return event
}

// startApp starts the Tview app with the layout
func (g *Game) Start(app *tview.Application, authToken []byte) {
	g.app = app
	g.app.Stop()
	g.gameServer.SetOnStateChange(func(gs i.GameState) {
		g.renderMaze(gs)
		g.renderScoreboard(gs)
		g.app.Draw()
	})
	g.gameServer.SetOnPingResult(func(ping int64) {
		g.renderPing(ping)
		g.app.Draw()
	})

	statusBar := tview.NewTextView().
		SetText("[#6c7086]↑↓←→ / hjkl: Move   Ctrl+C: Quit[#cdd6f4]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	statusBar.SetBackgroundColor(catBase)

	rightSidebar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(g.scoreTV, 0, 2, false).
		AddItem(g.pingTV, 0, 1, false)
	rightSidebar.SetBackgroundColor(catBase)

	topRow := tview.NewFlex().
		AddItem(g.mazeTV, 0, 3, true).
		AddItem(rightSidebar, 0, 1, false)
	topRow.SetBackgroundColor(catBase)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topRow, 0, 1, true).
		AddItem(statusBar, 1, 0, false)
	layout.SetBackgroundColor(catBase)

	g.app.SetInputCapture(g.handleInput)
	g.mazeTV.SetText("loading...")
	go func() {
		_ = g.gameServer.Start(authToken)
	}()

	go func() {
		if err := app.SetRoot(layout, true).Run(); err != nil {
			panic(err)
		}
	}()

	for range g.stopChan {
		g.app.Stop()
		return
	}
}

func (g *Game) renderScoreboard(gs i.GameState) {
	players := gs.RetrivePlayers()

	// Sort by score if not then by ID for consistency
	sort.Slice(players, func(i, j int) bool {
		return (players[i].GetReward() > players[j].GetReward() ||
			players[i].GetID().String() < players[j].GetID().String())
	})

	g.scoreTV.Clear()

	// Set headers
	g.scoreTV.SetCell(0, 0, tview.NewTableCell("Player").SetTextColor(catMauve).SetAlign(tview.AlignCenter))
	g.scoreTV.SetCell(0, 1, tview.NewTableCell("Score").SetTextColor(catMauve).SetAlign(tview.AlignCenter))

	g.initPlayerDisplays(players)

	// Add player scores
	for i, player := range players {
		label := g.playerDisplays[player.GetID()].scoreboard
		g.scoreTV.SetCell(i+1, 0, tview.NewTableCell(label).SetAlign(tview.AlignLeft))
		g.scoreTV.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("%d", player.GetReward())).SetTextColor(tcell.GetColor("#a6e3a1")).SetAlign(tview.AlignRight))
	}
}

func (g *Game) renderPing(ping int64) {
	text := fmt.Sprintf("[#cba6f7]PING\n\n[#cdd6f4]Ping: [#89dceb]%dms", ping)
	g.pingTV.SetText(text)
	g.app.Draw()
}

// renderMaze renders the maze into a string
func (g *Game) renderMaze(gs i.GameState) {
	var builder strings.Builder
	grid := mazeGridRepr(gs)
	playersRpr := g.playerMap(gs)

	// Top border using HasNorthWall from the first row of cells
	firstRow := gs.RetriveMaze().RetriveGrid()[0]
	for _, cell := range firstRow {
		if cell.HasWestWall() || cell.HasNorthWall() {
			builder.WriteString("[:#45475a:]  [:#11111b:]")
		} else {
			builder.WriteString("  ")
		}
		if cell.HasNorthWall() {
			builder.WriteString("[:#45475a:]  [:#11111b:]")
		} else {
			builder.WriteString("  ")
		}
	}
	builder.WriteString("[:#45475a:]  [:#11111b:]\n") // East border corner

	for y, r := range grid {
		for x, c := range r {
			if repr, ok := playersRpr[fmt.Sprintf("%d,%d", x, y)]; ok {
				builder.WriteString(repr) // Player position
			} else if c == -1 {
				builder.WriteString("[:#45475a:]  [:#11111b:]") // Wall (Surface1 bg)
			} else if c == 1 {
				builder.WriteString("[#cdd6f4] ●[#cdd6f4]") // Reward 1 (Text)
			} else if c == 5 {
				builder.WriteString("[#f9e2af] ●[#cdd6f4]") // Reward 5 (Yellow)
			} else {
				builder.WriteString("  ") // Empty space
			}
		}
		builder.WriteString("\n")
	}
	g.mazeTV.SetText(builder.String())
}

// initPlayerDisplays builds display strings for all players (called once on first render).
// Other-player indices are assigned independently of the current player's position in the list,
// so numbering is always contiguous (P1, P2, …) and color assignment is correct.
func (g *Game) initPlayerDisplays(players []i.Player) {
	if g.playerDisplays != nil {
		return
	}
	g.playerDisplays = make(map[uuid.UUID]playerDisplay)

	// Catppuccin Mocha accent palette — excludes green (reserved for self)
	colors := []string{"#cba6f7", "#89b4fa", "#f9e2af", "#fab387", "#f38ba8", "#74c7ec"}
	other := 0
	for _, p := range players {
		if p.GetID() == g.playerID {
			g.playerDisplays[p.GetID()] = playerDisplay{
				maze:       "[#a6e3a1::b] ◆[:#11111b:][-::-]", // bold green diamond — "you"
				scoreboard: "[#a6e3a1::b]You[-::-]",
			}
		} else {
			color := colors[other%len(colors)]
			num := other + 1
			g.playerDisplays[p.GetID()] = playerDisplay{
				maze:       fmt.Sprintf("[%s::b] %d[:#11111b:][-::-]", color, num),
				scoreboard: fmt.Sprintf("[%s::b]P%d[-::-]", color, num),
			}
			other++
		}
	}
}

func (g *Game) playerMap(gs i.GameState) map[string]string {
	players := gs.RetrivePlayers()
	g.initPlayerDisplays(players)
	rprMap := make(map[string]string)
	for _, p := range players {
		key := fmt.Sprintf("%d,%d", p.RetrivePos().GetCol()*2+1, p.RetrivePos().GetRow()*2)
		rprMap[key] = g.playerDisplays[p.GetID()].maze
	}
	return rprMap
}

// mazeGridRepr generates a grid representation from the maze skipping players.
func mazeGridRepr(gs i.GameState) [][]int {
	var grid [][]int
	for _, row := range gs.RetriveMaze().RetriveGrid() {
		r := make([]int, 0)
		for _, cell := range row {
			if cell.HasWestWall() {
				r = append(r, -1)
			} else {
				r = append(r, 0)
			}
			r = append(r, int(cell.GetReward()))
		}

		r = append(r, -1)
		grid = append(grid, r)
		r = make([]int, 0)
		for _, cell := range row {
			if cell.HasWestWall() || cell.HasSouthWall() {
				r = append(r, -1)
			} else {
				r = append(r, 0)
			}

			if cell.HasSouthWall() {
				r = append(r, -1)
			} else {
				r = append(r, 0)
			}
		}
		r = append(r, -1)
		grid = append(grid, r)
	}
	return grid
}
