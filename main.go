package main

import (
	"fmt"
	"md_to_pdf/converter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"os"
	"path/filepath"
	"time"
)

func main() {
	a := app.New()
	w := a.NewWindow("Markdown to PDF")

	var inputPath string
	var outputPath string

	inputLabel := widget.NewLabel("No file selected")
	outputLabel := widget.NewLabel("No output selected")
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	selectInput := widget.NewButton("Select Markdown", func() {
		d := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil || file == nil {
				return
			}
			inputPath = file.URI().Path()
			inputLabel.SetText(inputPath)
		}, w)
		uri, err := storage.ListerForURI(storage.NewFileURI(exeDir))
		if err == nil {
			d.SetLocation(uri)
		}

		d.Show()
	})

	selectOutput := widget.NewButton("Select Output Location", func() {
		dialog.NewFileSave(func(file fyne.URIWriteCloser, err error) {
			if err != nil || file == nil {
				return
			}
			outputPath = file.URI().Path()
			outputLabel.SetText(outputPath)
		}, w).Show()
	})

	convertBtn := widget.NewButton("Convert", func() {
		if inputPath == "" || outputPath == "" {
			dialog.ShowInformation("Error", "Please select input and output", w)
			return
		}

		if filepath.Ext(outputPath) != ".pdf" {
			outputPath += ".pdf"
		}

		start := time.Now()

		err := converter.Convert(inputPath, outputPath)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		elapsed := time.Since(start)

		dialog.ShowInformation(
			"Success",
			fmt.Sprintf("PDF created successfully!\nTime taken: %v", elapsed),
			w,
		)
	})

	content := container.NewVBox(
		selectInput,
		inputLabel,
		selectOutput,
		outputLabel,
		convertBtn,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 500))
	w.ShowAndRun()
}
