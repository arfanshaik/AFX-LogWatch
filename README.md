# ⚡ AFX LogWatch

**AFX LogWatch** is a fast, lightweight real-time log monitoring CLI written in **Go**. It follows one or more log files, detects common log levels, highlights important events, filters noise with regular expressions, exports JSON Lines, handles common log rotation cases, and can print live event statistics.

> **Monitor. Filter. Debug. Faster.**

## Why this project?

Large log files become difficult to read when errors, warnings and normal messages are mixed together. AFX LogWatch gives developers a focused terminal view of changing logs without requiring a large observability stack.

## Features

- Written in Go with the standard library only
- Follow multiple log files in real time
- Detect `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR` and `FATAL`
- Color-coded terminal output
- Include / exclude regex filters
- Filter by one or more log levels
- Read from the end of a file by default
- `--from-start` option for existing content
- JSON Lines output for scripts and pipelines
- Periodic event statistics
- Detect file replacement / truncation for common log rotation workflows
- Graceful shutdown with `Ctrl+C`
- Cross-platform core design
- Unit tests and GitHub Actions CI

## Project Structure

```text
AFX-LogWatch/
├── cmd/afx-logwatch/main.go
├── internal/
│   ├── formatter/formatter.go
│   ├── parser/parser.go
│   ├── stats/stats.go
│   └── watcher/watcher.go
├── examples/app.log
├── .github/workflows/go.yml
├── go.mod
├── LICENSE
└── README.md
```

## Requirements

- Go 1.22 or newer

## Run without installing

```bash
go run ./cmd/afx-logwatch watch --from-start examples/app.log
```

## Build

### Windows

```powershell
go build -o AFX-LogWatch.exe ./cmd/afx-logwatch
```

### Linux / macOS

```bash
go build -o afx-logwatch ./cmd/afx-logwatch
```

## Usage

```text
afx-logwatch watch [options] <file> [file...]
afx-logwatch scan  [options] <file>
afx-logwatch version
```

### Follow a log file

```bash
afx-logwatch watch app.log
```

### Read old lines first, then continue watching

```bash
afx-logwatch watch --from-start app.log
```

### Only errors and warnings

```bash
afx-logwatch watch --level error,warn app.log
```

### Search for important messages

```bash
afx-logwatch watch --include "timeout|failed|refused" app.log
```

### Hide noisy health-check messages

```bash
afx-logwatch watch --exclude "health|heartbeat" app.log
```

### Watch multiple logs

```bash
afx-logwatch watch api.log worker.log database.log
```

### JSON output

```bash
afx-logwatch watch --json app.log
```

### Periodic statistics

```bash
afx-logwatch watch --stats-every 30s app.log
```

### Scan an existing file once

```bash
afx-logwatch scan --level error examples/app.log
```

## Example output

```text
app.log            INFO      2026-09-21 INFO server started
app.log            WARN      2026-09-21 WARN request latency exceeded 900ms
app.log            ERROR     2026-09-21 ERROR timeout contacting payment service
```

## Development

Format, test and build the project:

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/afx-logwatch
```

## Good use cases

AFX LogWatch can help with local development, backend API debugging, server logs, application diagnostics, CI output, game/server logs and lightweight production troubleshooting where a full log platform is unnecessary.

## Roadmap

Possible future additions include structured JSON log parsing, configurable timestamp parsing, saved filter profiles, WebSocket streaming, a local web dashboard and optional desktop notifications.

## License

MIT License — see [LICENSE](LICENSE).

---

Built as part of the **AFX** developer-tool ecosystem.
