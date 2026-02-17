# md2pdf

A minimal CLI tool that converts Markdown files to PDF using Pandoc and XeLaTeX.

> Built while learning Go — focused on real-world tooling and understanding how things work under the hood.

---

## Requirements

- Go 1.20+
- [Pandoc](https://pandoc.org)
- LaTeX distribution (MiKTeX or TeX Live) with XeLaTeX
- A Unicode font (recommended: Noto Sans)

---

## Usage

```bash
# Run directly
go run main.go -in file.md

# Or build first
go build -o md2pdf
./md2pdf -in file.md
```

Internally runs:
```bash
pandoc input.md -o output.pdf --pdf-engine=xelatex -V geometry:margin=1in
```

---

## What I learned

- Running external commands with `exec.Command`
- How Pandoc converts Markdown → LaTeX → PDF
- Unicode and font handling in XeLaTeX
- Structuring a small but complete CLI tool

---

## Roadmap

**v1.x**
- Custom margin flag
- Verbose/debug mode
- Unit tests

**Next Goals** — Desktop app (Wails)
- File picker, drag & drop
- Output location selector
- Progress feedback and status log

---

## Author

Suraj — learning Go by building real tools.