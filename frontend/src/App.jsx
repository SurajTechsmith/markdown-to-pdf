import { useState } from "react"
import "./App.css"
import {
  SelectMarkdownFile,
  SavePDF,
  ConvertMarkdown,
  SelectFolder,
  ConvertFolder,
} from "../wailsjs/go/main/App"

function statusState(status) {
  if (status === "Ready") return "ready"
  if (status.includes("Converting")) return "working"
  if (status.includes("✓")) return "success"
  if (status.includes("✗")) return "error"
  return "ready"
}

function truncatePath(p) {
  if (!p) return ""
  const parts = p.replace(/\\/g, "/").split("/")
  return parts.length <= 3 ? p : "…/" + parts.slice(-2).join("/")
}

function parseError(err) {
  if (typeof err !== "string") return "Conversion failed ✗"
  if (err.includes("xelatex") || err.includes("pdf-engine")) return "Missing PDF engine (xelatex) ✗"
  if (err.includes("pandoc")) return "Pandoc error — check terminal ✗"
  return "Conversion failed ✗"
}

export default function App() {
  const [mode, setMode] = useState("file")
  const [inputFile, setInputFile] = useState("")
  const [outputFile, setOutputFile] = useState("")
  const [inputFolder, setInputFolder] = useState("")
  const [outputFolder, setOutputFolder] = useState("")
  const [status, setStatus] = useState("Ready")
  const [duration, setDuration] = useState("")
  const [loading, setLoading] = useState(false)
  const [folderResults, setFolderResults] = useState(null)

  function switchMode(m) {
    setMode(m)
    setStatus("Ready")
    setDuration("")
    setFolderResults(null)
  }

  async function convert() {
    setDuration("")

    if (mode === "file") {
      if (!inputFile || !outputFile) { setStatus("Select input and output file"); return }
      setLoading(true)
      setStatus("Converting…")
      try {
        const elapsed = await ConvertMarkdown(inputFile, outputFile)
        setStatus("PDF created successfully ✓")
        setDuration(elapsed)
      } catch (err) {
        setStatus(parseError(err))
      } finally {
        setLoading(false)
      }
    } else {
      if (!inputFolder) { setStatus("Select an input folder"); return }
      setLoading(true)
      setStatus("Converting folder…")
      setFolderResults(null)
      try {
        const results = await ConvertFolder(inputFolder, outputFolder)
        const failed = results.filter(r => r.Error).length
        const ok = results.length - failed
        const total = results.reduce((sum, r) => r.Duration ? sum + parseFloat(r.Duration) : sum, 0)
        setFolderResults(results)
        setStatus(failed === 0 ? `${ok} file${ok !== 1 ? "s" : ""} converted ✓` : `${ok} succeeded, ${failed} failed ✗`)
        setDuration(`${total.toFixed(2)}s total`)
      } catch (err) {
        setStatus(typeof err === "string" ? err + " ✗" : "Folder conversion failed ✗")
      } finally {
        setLoading(false)
      }
    }
  }

  const canConvert = mode === "file" ? !!(inputFile && outputFile) : !!inputFolder
  const state = statusState(status)

  return (
    <div className="app">
      <div className="card">

        <div className="card-header">
          <div className="logo">
            <span className="logo-title">Pandoc</span>
            <span className="logo-badge">MD → PDF</span>
          </div>
          <div className="subtitle">Document Converter</div>
        </div>

        <div className="tab-row">
          <button className={`tab ${mode === "file" ? "active" : ""}`} onClick={() => switchMode("file")}>Single File</button>
          <button className={`tab ${mode === "folder" ? "active" : ""}`} onClick={() => switchMode("folder")}>Folder</button>
        </div>

        <div className="card-body">
          {mode === "file" ? (
            <>
              <div className="field">
                <div className="field-label">Input</div>
                <div className="field-row" onClick={() => SelectMarkdownFile().then(f => f && setInputFile(f))}>
                  <span className="field-icon">◆</span>
                  <span className={`field-value ${inputFile ? "set" : ""}`}>{inputFile ? truncatePath(inputFile) : "Select markdown file…"}</span>
                </div>
              </div>
              <div className="field">
                <div className="field-label">Output</div>
                <div className="field-row" onClick={() => SavePDF().then(f => f && setOutputFile(f))}>
                  <span className="field-icon">◇</span>
                  <span className={`field-value ${outputFile ? "set" : ""}`}>{outputFile ? truncatePath(outputFile) : "Choose save location…"}</span>
                </div>
              </div>
            </>
          ) : (
            <>
              <div className="field">
                <div className="field-label">Input Folder</div>
                <div className="field-row" onClick={() => SelectFolder().then(f => f && setInputFolder(f))}>
                  <span className="field-icon">▣</span>
                  <span className={`field-value ${inputFolder ? "set" : ""}`}>{inputFolder ? truncatePath(inputFolder) : "Select folder with .md files…"}</span>
                </div>
              </div>
              <div className="field">
                <div className="field-label">Output Folder</div>
                <div className="field-row" onClick={() => SelectFolder().then(f => f && setOutputFolder(f))}>
                  <span className="field-icon">▢</span>
                  <span className={`field-value ${outputFolder ? "set" : ""}`}>{outputFolder ? truncatePath(outputFolder) : "Same as input folder (default)"}</span>
                </div>
                <div className="field-hint">Leave empty to save PDFs alongside source files</div>
              </div>
            </>
          )}

          <div className="divider" />

          <button
            className={`convert-btn ${loading ? "loading" : ""}`}
            onClick={convert}
            disabled={!canConvert || loading}
          >
            {loading ? "Converting…" : "Convert to PDF"}
          </button>

          {folderResults && (
            <div className="results">
              <div className="results-header">
                <span>{folderResults.length} files processed</span>
                {duration && <span className="results-header-duration">{duration}</span>}
              </div>
              {folderResults.map((r, i) => (
                <div className="result-row" key={i}>
                  <span className="result-name">{r.File}</span>
                  {r.Duration && <span className="result-duration">{r.Duration}</span>}
                  <span className={`result-badge ${r.Error ? "fail" : "ok"}`}>{r.Error ? "failed" : "ok"}</span>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="status-bar">
          <div className={`status-dot ${state}`} />
          <span className="status-text">{status}</span>
          {duration && mode === "file" && <span className="status-duration">{duration}</span>}
        </div>

      </div>
    </div>
  )
}