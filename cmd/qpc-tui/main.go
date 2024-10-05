package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	tea "github.com/charmbracelet/bubbletea"

	"qpc-tui/internal/app"
)

const (
	host = "0.0.0.0"
	port = "22"
)

func main() {
	// Initialize the server
	s, err := wish.NewServer(
		// Set the address to the host and port, using net.JoinHostPort to combine them
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			// Initialize the Bubble Tea middleware with a custom function that initializes the Bubble Tea model and options
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				m, opts := app.InitialModel(s)
				return m, opts
			}),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	// Start the server in a separate goroutine, log any errors and notify the main goroutine when done
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		// ListenAndServe is a blocking call that starts the server and returns an error if it occurs
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	// Wait for the server to be stopped
	<-done
	log.Info("Stopping SSH server")
	// Create a context with a timeout of 30 seconds to stop the server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	// Shutdown the server with the context
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}