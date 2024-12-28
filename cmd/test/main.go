package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	fmt.Println("Starting test app")

	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(hello)

	fmt.Println("Showing window")
	w.Show()

	fmt.Println("Running app")
	a.Run()
}
