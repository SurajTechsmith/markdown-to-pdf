
# Markdown to PDF Converter (v1.0.0)

A simple CLI tool built in Go that converts Markdown files into PDF using Pandoc and XeLaTeX.

> Built while learning Go — focused on real-world tooling, system interaction, and understanding how tools work under the hood.

---

## 🚀 Version

Current version: **v1.0.0**  
This is the first stable CLI release.

This version focuses on correctness, Unicode handling, and reliable PDF generation.

---

## ✨ Features

- Convert `.md` files to `.pdf`
- Uses `xelatex` for full Unicode support
- Custom margin configuration
- Handles Pandoc execution errors properly
- Clean and minimal CLI interface

---

## 🛠 Requirements

Make sure the following are installed:

- Go (1.20+ recommended)
- Pandoc
- MiKTeX (or any LaTeX distribution)
- XeLaTeX engine enabled
- A Unicode font (recommended: Noto Sans)

---

## 📦 Installation

Clone the repository:

```bash
git clone https://github.com/SurajTechsmith/markdown-to-pdf
cd markdown-to-pdf
````

---

## ▶️ Usage

Run directly:

```bash
go run main.go -in file.md
```

Or build an executable:

```bash
go build -o md2pdf
./md2pdf -in file.md
```

The tool internally executes:

```bash
pandoc input.md -o output.pdf --pdf-engine=xelatex -V geometry:margin=1in
```

---

## 🧠 What This Project Helped Me Learn

* Executing external commands using Go (`exec.Command`)
* Understanding how Pandoc converts Markdown → LaTeX → PDF
* Unicode handling in XeLaTeX
* Debugging CLI exit errors
* Font configuration issues in MiKTeX
* Structuring a small but complete CLI tool

---

## 🔮 Future Goals (Roadmap)

### v1.x Improvements

* Add custom margin flag from CLI
* Improve error output formatting
* Add verbose/debug mode
* Add unit tests for argument parsing

### v2.0 – Desktop Application

Convert this CLI tool into a cross-platform desktop application using Fyne:

* File picker for selecting Markdown files
* Output location selector
* Convert button with progress feedback
* Status/log panel
* Font selector option
* Drag & drop support
* Cross-platform builds (Windows / Linux / macOS)

### Long-Term Vision

* Batch file conversion
* Configuration file support
* Predefined themes/templates
* GitHub release binaries
* Possibly turn it into a lightweight Markdown publishing tool

---

## 👤 Author

Suraj Yadav
Learning Go by building real tools 🚀

