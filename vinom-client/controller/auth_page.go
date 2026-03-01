package controller

import (
	"github.com/beka-birhanu/vinom-client/dmn"
	"github.com/beka-birhanu/vinom-client/service/i"
	"github.com/rivo/tview"
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
	header := tview.NewTextView().SetText("Login / Sign Up").SetTextAlign(tview.AlignCenter)
	footer := tview.NewTextView().SetText("").SetTextAlign(tview.AlignLeft)

	form := tview.NewForm()
	form.AddInputField("Username:", "", 20, nil, nil)
	form.AddPasswordField("Password:", "", 20, '*', nil)

	form.AddButton("Login", func() {
		username := form.GetFormItem(0).(*tview.InputField).GetText()
		password := form.GetFormItem(1).(*tview.InputField).GetText()
		player, token, err := a.authService.Login(username, password)
		if err != nil {
			footer.SetText(err.Error())
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

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(footer, 0, 1, false)

	return flex
}

func (a *AuthPage) signUpForm(app *tview.Application) tview.Primitive {
	header := tview.NewTextView().SetText("Sign Up").SetTextAlign(tview.AlignCenter)
	footer := tview.NewTextView().SetText("").SetTextAlign(tview.AlignLeft)

	form := tview.NewForm()
	form.AddInputField("Username:", "", 20, nil, nil)
	form.AddPasswordField("Password:", "", 20, '*', nil)

	form.AddButton("Register", func() {
		username := form.GetFormItem(0).(*tview.InputField).GetText()
		password := form.GetFormItem(1).(*tview.InputField).GetText()
		err := a.authService.Register(username, password)
		if err != nil {
			footer.SetText(err.Error())
			return
		}
		app.SetRoot(a.signInForm(app), true)
	})

	form.AddButton("Back", func() {
		app.SetRoot(a.signInForm(app), true)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(footer, 0, 1, false)

	return flex
}
