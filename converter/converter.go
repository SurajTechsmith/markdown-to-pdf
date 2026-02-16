package converter

import (
	"fmt"
	"log"
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

	cmd := exec.Command("pandoc", inputPath, "-o", outputPath, "--pdf-engine=xelatex", "-V geometry:margin=0.7in")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Pandoc error:")
		fmt.Println(string(output))
		log.Fatal(err)
	}

	return nil
}
