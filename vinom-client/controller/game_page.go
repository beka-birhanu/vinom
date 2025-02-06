package controller

import (
	"fmt"
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
	tcell.KeyLeft:  "East",
	tcell.KeyRight: "West",
}

// Additional handling for Vim motions
var vimDirections = map[rune]string{
	'k': "North", // Vim motion: k for up
	'j': "South", // Vim motion: j for down
	'h': "West",  // Vim motion: h for left
	'l': "East",  // Vim motion: l for right
}

// Game holds the maze, score, ping, and player information
type Game struct {
	gameServer   i.GameServer
	playerColors map[uuid.UUID]string
	playerID     uuid.UUID
	mazeTV       *tview.TextView
	scoreTV      *tview.TextView
	pingTV       *tview.TextView
	stopChan     chan struct{}
}

// NewGame creates a new MazeGame instance
func NewGame(gmSrvr i.GameServer, pID uuid.UUID) *Game {
	return &Game{
		gameServer: gmSrvr,
		playerID:   pID,
		mazeTV:     tview.NewTextView().SetDynamicColors(true),
		scoreTV:    tview.NewTextView().SetDynamicColors(true),
		pingTV:     tview.NewTextView().SetDynamicColors(true),
		stopChan:   make(chan struct{}),
	}
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
func (g *Game) StartApp(app *tview.Application) {
	// Combine maze, scoreboard, and ping into a Flex layout
	layout := tview.NewFlex().
		AddItem(g.mazeTV, 0, 3, true).   // Maze occupies 3/4 of the screen width
		AddItem(g.scoreTV, 0, 1, false). // Scoreboard
		AddItem(g.pingTV, 0, 1, false)   // Ping

	app.SetInputCapture(g.handleInput)
	g.mazeTV.SetText("loading...")

	go func() {
		if err := app.SetRoot(layout, true).Run(); err != nil {
			panic(err)
		}
	}()
	g.listenAndRender()

	for range g.stopChan {
		app.Stop()
		return
	}
}

func (g *Game) listenAndRender() {
	go func() {
		for gameState := range g.gameServer.StateChan() {
			g.renderMaze(gameState)
			g.renderScoreboard(gameState)
		}
	}()

	go func() {
		for ping := range g.gameServer.PingChan() {
			g.renderPing(ping)
		}
	}()
}

func isPlayerPos(x int, y int, players []i.Player) i.Player {
	for _, player := range players {
		if int(player.RetrivePos().GetRow()) == x && int(player.RetrivePos().GetCol()) == y {
			return player
		}
	}

	return nil
}

func (g *Game) playerRepr(pID uuid.UUID, players []i.Player) string {
	if g.playerColors == nil {
		g.playerColors = make(map[uuid.UUID]string)
		colors := [6]string{"yellow", "orange", "lime", "purple", "magenta"}
		i := 0
		for _, p := range players {
			if p.GetID() == g.playerID {
				g.playerColors[p.GetID()] = "⭕"
			} else {
				g.playerColors[p.GetID()] = fmt.Sprintf("[%s]P%d[black]", colors[i], i+1)
			}
			i++
		}
	}

	return g.playerColors[pID]
}

// gridRepr generates a grid representation from the maze
func gridRepr(gs i.GameState) [][]int {
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

// renderMaze renders the maze into a string
func (g *Game) renderMaze(gs i.GameState) {
	var builder strings.Builder
	// Top border
	grid := gridRepr(gs)
	players := gs.RetrivePlayers()

	builder.WriteString(strings.Repeat("[:blue]  [:black]", len(grid[0])) + "\n")
	for y, r := range grid {
		for x, c := range r {
			if player := isPlayerPos(x, y, players); player != nil {
				builder.WriteString(g.playerRepr(player.GetID(), players)) // Player position
			} else if c == -1 {
				builder.WriteString("[:blue]  [:black]") // Wall
			} else if c == 1 {
				builder.WriteString("[white] ●[black]") // Reward 1
			} else if c == 5 {
				builder.WriteString("[yellow] ●[black]") // Reward 5
			} else {
				builder.WriteString("  ") // Empty space
			}
		}
		builder.WriteString("\n")
	}
	g.mazeTV.SetText(builder.String())
}

func (mg *Game) renderScoreboard(gs i.GameState) {
	text := fmt.Sprintf("[yellow]SCOREBOARD\n\n[white]Score: [green]%d", gs.RetrivePlayers()[0].GetReward())
	mg.scoreTV.SetText(text)
}

func (mg *Game) renderPing(ping int64) {
	text := fmt.Sprintf("[yellow]PING\n\n[white]Ping: [cyan]%dms", ping)
	mg.pingTV.SetText(text)
}
