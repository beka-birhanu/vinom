package controller

import "github.com/rivo/tview"

type PageController interface {
	Start(*tview.Application) error
}
