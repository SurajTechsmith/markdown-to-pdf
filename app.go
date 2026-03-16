package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"markdown2pdf/converter"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectMarkdownFile() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Markdown File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Markdown (*.md, *.markdown)", Pattern: "*.md;*.markdown"},
		},
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func (a *App) SavePDF() (string, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save PDF As",
		DefaultFilename: "output.pdf",
		Filters: []runtime.FileFilter{
			{DisplayName: "PDF (*.pdf)", Pattern: "*.pdf"},
		},
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func (a *App) SelectFolder() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Folder",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func (a *App) ConvertMarkdown(inputPath, outputPath string) (string, error) {
	elapsed, err := converter.Convert(inputPath, outputPath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2fs", elapsed.Seconds()), nil
}

func (a *App) ConvertFolder(inputDir, outputDir string) ([]converter.FolderResult, error) {
	return converter.ConvertFolder(inputDir, outputDir)
}
