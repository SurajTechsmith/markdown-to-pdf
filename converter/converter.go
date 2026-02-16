package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Convert(inputPath string, outputPath string) error {
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		base := inputPath[:len(inputPath)-len(ext)]
		outputPath = base + ".pdf"
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("File does not exist")

	}

	cmd := exec.Command("pandoc", inputPath, "-o", outputPath, " --pdf-engine=xelatex", "-V geometry:margin=0.7in")

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run pandoc: %w", err)
	}

	return nil
}
