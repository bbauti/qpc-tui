# qpc-tui

A Terminal User Interface (TUI) application built with Bubble Tea, accessible via SSH.

## Prerequisites

- Go 1.23.0 or later

## Installation

1. **Clone the repository** (if you haven't already):
```bash
git clone <repository-url>
cd qpc-tui
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Generate SSH host keys** (required for the SSH server):
```bash
mkdir -p .ssh && ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""
```

## Running the Application

### Option 1: Run directly with Go

```bash
go run cmd/qpc-tui/main.go
```

### Option 2: Build and run

```bash
go build -o qpc-tui cmd/qpc-tui/main.go
./qpc-tui
```

## Connecting to the Application

The server runs on `localhost:22` by default. Connect via SSH from another terminal:

```bash
ssh localhost -p 22
```

### Running on a different port

Port 22 requires root/sudo privileges. To run without elevated privileges:

1. Edit `cmd/qpc-tui/main.go` and change line 26:
```go
port = "2222"  // or any port > 1024
```

2. Connect using:
```bash
ssh localhost -p 2222
```

## Project Structure

```
.
├── cmd/
│   └── qpc-tui/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/                 # Application logic
│   │   ├── delegate.go
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   ├── scraper/             # Web scraping functionality
│   │   └── scraper.go
│   └── ui/                  # UI components
│       └── keys.go
├── go.mod
└── go.sum
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Wish](https://github.com/charmbracelet/wish) - SSH server
- [Colly](https://github.com/gocolly/colly) - Web scraping
