package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	appName = "Fynancial App"

	appHSize = 275
	appVSize = 250
)

var (
	balance = 0.0

	colorGreen = color.RGBA{R: 0, G: 255, B: 0, A: 100}

	colorRed = color.RGBA{R: 255, G: 0, B: 0, A: 100}

	colorGreenBlue = color.RGBA{R: 0, G: 255, B: 255, A: 100}

	errDialogOnParse = errors.New("invalid entry: only unnsigned numbers")
)

func main() {
	a := app.New()

	// Window
	w := setWindow(a)

	// Balance area
	balanceText := setBalanceText()

	// Withdraw area
	withdrawVBox := setWithdrawArea(w, balanceText)

	// Deposit area
	depositHVox := setDepositArea(w, balanceText)

	// Main Container
	mainContainer := setMainContainer(3, balanceText, withdrawVBox, depositHVox)

	w.SetContent(mainContainer)
	w.ShowAndRun()
}

func setWindow(a fyne.App) fyne.Window {
	w := a.NewWindow(appName)
	w.Resize(fyne.NewSize(appHSize, appVSize))
	iconResource, err := fyne.LoadResourceFromPath("./money-bag.png")
	if err != nil {
		log.Fatal(err)
	}
	w.SetIcon(iconResource)

	return w
}

func setBalanceText() *canvas.Text {
	balanceToString := fmt.Sprintf("R$ %s", strconv.FormatFloat(balance, 'f', 2, 64))
	balanceText := canvas.NewText(balanceToString, colorGreenBlue)
	balanceText.Alignment = fyne.TextAlignCenter
	balanceText.TextStyle = fyne.TextStyle{Bold: true}
	balanceText.TextSize = 32

	return balanceText
}

func updateBalanceText(balance float64, balanceText *canvas.Text) {
	balanceToString := fmt.Sprintf("R$ %s", strconv.FormatFloat(balance, 'f', 2, 64))
	switch {
	case balance > 0.0:
		balanceText.Text = balanceToString
		balanceText.Color = colorGreen
	case balance < 0.0:
		balanceText.Text = balanceToString
		balanceText.Color = colorRed
	default:
		balanceText.Text = balanceToString
		balanceText.Color = colorGreenBlue
	}

	balanceText.Refresh()
}

func setWithdrawArea(w fyne.Window, balanceText *canvas.Text) *fyne.Container {
	withdrawEntry := widget.NewEntry()
	withdrawEntry.PlaceHolder = "amount to withdraw"
	withdrawButton := widget.NewButton("Withdraw", func() {
		value, err := strconv.ParseFloat(withdrawEntry.Text, 64)
		if err != nil {
			dialog.NewError(errDialogOnParse, w).Show()
			return
		}
		balance -= value
		updateBalanceText(balance, balanceText)
		withdrawEntry.SetText("")
	})
	return container.NewVBox(withdrawEntry, withdrawButton)
}

func setDepositArea(w fyne.Window, balanceText *canvas.Text) *fyne.Container {
	depositEntry := widget.NewEntry()
	depositEntry.PlaceHolder = "amount to deposit"
	depositButton := widget.NewButton("Deposit", func() {
		value, err := strconv.ParseFloat(depositEntry.Text, 64)
		if err != nil {
			dialog.NewError(errDialogOnParse, w).Show()
			return
		}
		balance += value
		updateBalanceText(balance, balanceText)
		depositEntry.SetText("")
	})
	return container.NewVBox(depositEntry, depositButton)
}

func setMainContainer(rows int, balanceText *canvas.Text, withdrawVBox, depositVBox *fyne.Container) *fyne.Container {
	return container.New(layout.NewGridLayoutWithRows(rows),
		balanceText,
		withdrawVBox,
		depositVBox,
	)
}
