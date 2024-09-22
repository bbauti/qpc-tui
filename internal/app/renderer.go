package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

type Renderer struct {
	ssh.Session
}

func (r Renderer) NewStyle() lipgloss.Style {
	return bubbletea.MakeRenderer(r.Session).NewStyle()
}