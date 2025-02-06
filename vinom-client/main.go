package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	signInPage := signInForm(app)

	if err := app.SetRoot(signInPage, true).Run(); err != nil {
		panic(err)
	}
}

func signInForm(app *tview.Application) tview.Primitive {
	header := tview.NewTextView().SetText("Login / Sign Up").SetTextAlign(tview.AlignCenter)

	form := tview.NewForm()
	form.AddInputField("Username:", "", 20, nil, nil)
	form.AddPasswordField("Password:", "", 20, '*', nil)

	form.AddButton("Login", func() {
		username := form.GetFormItem(0).(*tview.InputField).GetText()
		password := form.GetFormItem(1).(*tview.InputField).GetText()
		fmt.Printf("Logging in with %s and %s\n", username, password)
	})

	form.AddButton("Sign Up", func() {
		app.SetRoot(signUpForm(app), true)
	})

	form.AddButton("Quit", func() {
		app.Stop()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(form, 0, 1, true)

	return flex
}

func signUpForm(app *tview.Application) tview.Primitive {
	header := tview.NewTextView().SetText("Sign Up").SetTextAlign(tview.AlignCenter)

	form := tview.NewForm()
	form.AddInputField("Username:", "", 20, nil, nil)
	form.AddPasswordField("Password:", "", 20, '*', nil)

	form.AddButton("Register", func() {
		username := form.GetFormItem(0).(*tview.InputField).GetText()
		password := form.GetFormItem(1).(*tview.InputField).GetText()
		fmt.Printf("Registering %s with password %s\n", username, password)
	})

	form.AddButton("Back", func() {
		app.SetRoot(signInForm(app), true)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(form, 0, 1, true)

	return flex
}
