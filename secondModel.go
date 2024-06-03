package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hex/data_model"
)

type secondModel struct {
	filteredServicesList []*struct {
		Service    data_model.Service
		IsSelected bool
		Input      string
	}
	pointer             int
	snapshotDataTextBox textinput.Model
}

func (m secondModel) Init() tea.Cmd {
	return nil
}

func (m secondModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else if msg.Type == tea.KeyEnter {
			if m.pointer >= len(m.filteredServicesList) {
				ti := textinput.New()
				ti.Focus()
				ti.Placeholder = "Enter Environment Name"
				return thirdModel{servicesList: m.filteredServicesList, envNameTextBox: ti}, nil
			} else {
				m.filteredServicesList[m.pointer].Input = m.snapshotDataTextBox.Value()
				m.snapshotDataTextBox.Reset()
				m.pointer++
			}
		} else {
			m.snapshotDataTextBox, cmd = m.snapshotDataTextBox.Update(msg)
		}
	}
	return m, cmd
}

func (m secondModel) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7FFFD4")).Bold(true)
	var styleMagenta = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	s := ""
	for _, service := range m.filteredServicesList {
		if len(service.Input) > 0 {
			s += fmt.Sprintf("%s: %s\n", styleMagenta.Render(service.Service.Name), service.Input)
		}
	}
	if m.pointer >= len(m.filteredServicesList) {
		s += "\nAll inputs complete, Press Enter to create deploy.json"
		return s
	}
	s += fmt.Sprintf("Input data for %s: %s", style.Render(m.filteredServicesList[m.pointer].Service.Name), m.snapshotDataTextBox.View())
	return s
}
