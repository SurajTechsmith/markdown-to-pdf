package converter

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed pandoc-bin/windows/pandoc.exe
//go:embed pandoc-bin/linux/pandoc
var pandocBin embed.FS

//go:embed pandoc-bin/windows/wkhtmltox/bin/wkhtmltopdf.exe
//go:embed pandoc-bin/linux/usr/local/bin/wkhtmltopdf
var wkhtmltopdf embed.FS

//go:embed pandoc-bin/style.css
var styleBin embed.FS

var (
	pandocPath string
	pandocOnce sync.Once
	pandocErr  error
)

var (
	wkhtmltopdfPath string
	wkhtmltopdfOnce sync.Once
	wkhtmltopdfErr  error
)
var (
	cssPath string
	cssOnce sync.Once
	cssErr  error
)

func getCSSPath() (string, error) {
	cssOnce.Do(func() {
		data, err := styleBin.ReadFile("pandoc-bin/style.css")
		if err != nil {
			cssErr = fmt.Errorf("embedded css not found: %w", err)
			return
		}

		tmpDir, err := os.MkdirTemp("", "css-wails-*")
		if err != nil {
			cssErr = fmt.Errorf("cannot create temp dir: %w", err)
			return
		}

		targetPath := filepath.Join(tmpDir, "style.css")
		if err := os.WriteFile(targetPath, data, 0o644); err != nil {
			cssErr = fmt.Errorf("cannot write css: %w", err)
			return
		}

		cssPath = targetPath
	})

	return cssPath, cssErr
}

func getWkhtmltopdfPath() (string, error) {
	wkhtmltopdfOnce.Do(func() {

		tmpDir, err := os.MkdirTemp("", "wkhtmltopdf-wails-*")
		if err != nil {

			wkhtmltopdfErr = fmt.Errorf("cannot create temp directory for wkhtmltopdf: %w", err)
			return
		}

		var embeddedFile, targetFile string

		switch runtime.GOOS {
		case "windows":
			embeddedFile = "pandoc-bin/windows/wkhtmltox/bin/wkhtmltopdf.exe"
			targetFile = "wkhtmltopdf.exe"
		case "linux":
			embeddedFile = "pandoc-bin/linux/usr/local/bin/wkhtmltopdf"
			targetFile = "wkhtmltopdf"
		default:

			wkhtmltopdfErr = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
			return
		}

		data, err := wkhtmltopdf.ReadFile(embeddedFile)
		if err != nil {
			wkhtmltopdfErr = fmt.Errorf("embedded wkhtmltopdf not found (%s): %w", embeddedFile, err)
			return
		}

		targetPath := filepath.Join(tmpDir, targetFile)
		if err := os.WriteFile(targetPath, data, 0o755); err != nil {
			wkhtmltopdfErr = fmt.Errorf("cannot write wkhtmltopdf to temp: %w", err)
			return
		}

		if runtime.GOOS != "windows" {
			if err := os.Chmod(targetPath, 0o755); err != nil {
				wkhtmltopdfErr = fmt.Errorf("cannot make wkhtmltopdf executable: %w", err)
				return
			}
		}

		wkhtmltopdfPath = targetPath
	})

	if wkhtmltopdfErr != nil {
		return "", wkhtmltopdfErr
	}

	return wkhtmltopdfPath, nil
}
func getPandocPath() (string, error) {
	pandocOnce.Do(func() {
		tmpDir, err := os.MkdirTemp("", "pandoc-wails-*")
		if err != nil {
			pandocErr = fmt.Errorf("cannot create temp directory for pandoc: %w", err)
			return
		}

		var embeddedFile, targetFile string

		switch runtime.GOOS {
		case "windows":
			embeddedFile = "pandoc-bin/windows/pandoc.exe"
			targetFile = "pandoc.exe"
		case "darwin":
			embeddedFile = "pandoc-bin/darwin/pandoc"
			targetFile = "pandoc"
		case "linux":
			embeddedFile = "pandoc-bin/linux/pandoc"
			targetFile = "pandoc"
		default:
			pandocErr = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
			return
		}

		data, err := pandocBin.ReadFile(embeddedFile)
		if err != nil {
			pandocErr = fmt.Errorf("embedded pandoc file not found (%s): %w", embeddedFile, err)
			return
		}

		targetPath := filepath.Join(tmpDir, targetFile)
		if err := os.WriteFile(targetPath, data, 0o755); err != nil {
			pandocErr = fmt.Errorf("cannot write pandoc binary to temp: %w", err)
			return
		}

		if runtime.GOOS != "windows" {
			if err := os.Chmod(targetPath, 0o755); err != nil {
				pandocErr = fmt.Errorf("cannot make pandoc executable: %w", err)
				return
			}
		}

		pandocPath = targetPath
	})

	if pandocErr != nil {
		return "", pandocErr
	}

	return pandocPath, nil
}

func Convert(inputPath, outputPath string) (time.Duration, error) {
	pandocBinPath, err := getPandocPath()
	if err != nil {
		return 0, fmt.Errorf("pandoc setup failed: %w", err)
	}

	wkhtmltopdfBinPath, err := getWkhtmltopdfPath()
	if err != nil {
		return 0, fmt.Errorf("wkhtmltopdf error: %w", err)
	}
	cssFilePath, err := getCSSPath()
	if err != nil {
		return 0, fmt.Errorf("css error: %w", err)
	}

	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = strings.TrimSuffix(inputPath, ext) + ".pdf"
	}

	inputPath, err = filepath.Abs(inputPath)
	if err != nil {
		return 0, fmt.Errorf("invalid input path: %w", err)
	}

	outputPath, err = filepath.Abs(outputPath)
	if err != nil {
		return 0, fmt.Errorf("invalid output path: %w", err)
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("input file does not exist: %s", inputPath)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return 0, fmt.Errorf("cannot create output directory: %w", err)
	}

	start := time.Now()
	cmd := exec.Command(
		pandocBinPath,
		inputPath,
		"-o", outputPath,
		"--pdf-engine="+wkhtmltopdfBinPath,
		"--css="+cssFilePath,
		"--pdf-engine-opt=--enable-local-file-access",
		"--self-contained",
		"--pdf-engine-opt=--margin-top", "--pdf-engine-opt=0.5in",
		"--pdf-engine-opt=--margin-bottom", "--pdf-engine-opt=0.5in",
		"--pdf-engine-opt=--margin-left", "--pdf-engine-opt=0.5in",
		"--pdf-engine-opt=--margin-right", "--pdf-engine-opt=0.5in",
	)

	cmd.SysProcAttr = getSysAttr()

	output, err := cmd.CombinedOutput()
	elapsed := time.Since(start)

	if err != nil {
		return 0, fmt.Errorf(
			"pandoc conversion failed (exit %d):\nCommand: %s\nOutput:\n%s\nError: %w",
			cmd.ProcessState.ExitCode(),
			strings.Join(cmd.Args, " "),
			string(output),
			err,
		)
	}

	return elapsed, nil
}

type FolderResult struct {
	File     string
	Error    string
	Duration string
}

func ConvertFolder(inputDir, outputDir string) ([]FolderResult, error) {
	inputDir, err := filepath.Abs(inputDir)
	if err != nil {
		return nil, fmt.Errorf("invalid input directory: %w", err)
	}

	if outputDir == "" {
		outputDir = inputDir
	} else {
		outputDir, err = filepath.Abs(outputDir)
		if err != nil {
			return nil, fmt.Errorf("invalid output directory: %w", err)
		}
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create output directory: %w", err)
	}

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read input directory: %w", err)
	}

	var results []FolderResult

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".md" && ext != ".markdown" {
			continue
		}

		inputPath := filepath.Join(inputDir, name)
		outputPath := filepath.Join(outputDir, strings.TrimSuffix(name, filepath.Ext(name))+".pdf")

		result := FolderResult{File: name}
		elapsed, err := Convert(inputPath, outputPath)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Duration = fmt.Sprintf("%.2fs", elapsed.Seconds())
		}
		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no markdown files found in: %s", inputDir)
	}

	return results, nil
}
