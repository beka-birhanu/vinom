package controller

import (
	"github.com/beka-birhanu/vinom-client/dmn"
	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Catppuccin Mocha
var (
	catBase    = tcell.GetColor("#1e1e2e")
	catSurface = tcell.GetColor("#313244")
	catBlue    = tcell.GetColor("#89b4fa")
	catLavend  = tcell.GetColor("#b4befe")
	catMauve   = tcell.GetColor("#cba6f7")
	catText    = tcell.GetColor("#cdd6f4")
)

type loginResponseHandler func(*dmn.Player, string)

type AuthPage struct {
	authService i.AuthServer
	onLogin     loginResponseHandler
}

func NewAuthPage(as i.AuthServer, onLogin loginResponseHandler) (*AuthPage, error) {
	return &AuthPage{
		authService: as,
		onLogin:     onLogin,
	}, nil
}

func (a *AuthPage) Start(app *tview.Application) error {
	if err := app.SetRoot(a.signInForm(app), true).Run(); err != nil {
		return err
	}
	return nil
}

func (a *AuthPage) signInForm(app *tview.Application) tview.Primitive {
	title := tview.NewTextView().
		SetText("[#cba6f7::b]V I N O M[-::-]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	title.SetBackgroundColor(catBase)

	footer := tview.NewTextView().
		SetText("").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	footer.SetBackgroundColor(catBase)

	form := tview.NewForm()
	form.SetBackgroundColor(catBase)
	form.SetFieldBackgroundColor(catSurface)
	form.SetFieldTextColor(catText)
	form.SetButtonBackgroundColor(catBlue)
	form.SetButtonTextColor(catBase)
	form.SetLabelColor(catLavend)
	form.AddInputField("Username:", "", 20, nil, nil)
	form.AddPasswordField("Password:", "", 20, '*', nil)

	form.AddButton("Login", func() {
		username := form.GetFormItem(0).(*tview.InputField).GetText()
		password := form.GetFormItem(1).(*tview.InputField).GetText()
		player, token, err := a.authService.Login(username, password)
		if err != nil {
			footer.SetText("[#f38ba8]" + err.Error() + "[#cdd6f4]")
			return
		}
		a.onLogin(player, token)
	})

	form.AddButton("Sign Up", func() {
		app.SetRoot(a.signUpForm(app), true)
	})

	form.AddButton("Quit", func() {
		app.Stop()
	})

	frame := tview.NewFrame(form).SetBorders(1, 1, 0, 0, 1, 1)
	frame.SetBorder(true)
	frame.SetTitle(" Login ")
	frame.SetBorderColor(catMauve)
	frame.SetTitleColor(catBlue)
	frame.SetBackgroundColor(catBase)

	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 3, 0, false).
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

func (a *AuthPage) signUpForm(app *tview.Application) tview.Primitive {
	title := tview.NewTextView().
		SetText("[#cba6f7::b]V I N O M[-::-]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	title.SetBackgroundColor(catBase)

	footer := tview.NewTextView().
		SetText("").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	footer.SetBackgroundColor(catBase)

	form := tview.NewForm()
	form.SetBackgroundColor(catBase)
	form.SetFieldBackgroundColor(catSurface)
	form.SetFieldTextColor(catText)
	form.SetButtonBackgroundColor(catBlue)
	form.SetButtonTextColor(catBase)
	form.SetLabelColor(catLavend)
	form.AddInputField("Username:", "", 20, nil, nil)
	form.AddPasswordField("Password:", "", 20, '*', nil)

	form.AddButton("Register", func() {
		username := form.GetFormItem(0).(*tview.InputField).GetText()
		password := form.GetFormItem(1).(*tview.InputField).GetText()
		err := a.authService.Register(username, password)
		if err != nil {
			footer.SetText("[#f38ba8]" + err.Error() + "[#cdd6f4]")
			return
		}
		app.SetRoot(a.signInForm(app), true)
	})

	form.AddButton("Back", func() {
		app.SetRoot(a.signInForm(app), true)
	})

	frame := tview.NewFrame(form).SetBorders(1, 1, 0, 0, 1, 1)
	frame.SetBorder(true)
	frame.SetTitle(" Sign Up ")
	frame.SetBorderColor(catMauve)
	frame.SetTitleColor(catBlue)
	frame.SetBackgroundColor(catBase)

	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 3, 0, false).
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
