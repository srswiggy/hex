package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type secondModel struct {
	filteredServicesList []*Service
	pointer              int
	snapshotDataTextBox  textinput.Model
}

func (m secondModel) Init() tea.Cmd {
	return nil
}

func initialSecondModel() secondModel {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "Enter Snapshot Tag for Service"

	servicesList := []*Service{
		{name: "finance-calcy-service", selected: true},
		{name: "finance-job-service", selected: true},
		{name: "finance-orchestrator", selected: true},
		{name: "finance-dashboard", selected: false},
	}

	var selectedServices []*Service
	for _, service := range servicesList {
		if service.selected {
			selectedServices = append(selectedServices, service)
		}
	}
	return secondModel{filteredServicesList: selectedServices, snapshotDataTextBox: ti}
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
				m.filteredServicesList[m.pointer].input = m.snapshotDataTextBox.Value()
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
		if len(service.input) > 0 {
			s += fmt.Sprintf("%s: %s\n", styleMagenta.Render(service.name), service.input)
		}
	}
	if m.pointer >= len(m.filteredServicesList) {
		s += "\nAll inputs complete, Press Enter to create deploy.json"
		return s
	}
	s += fmt.Sprintf("Input data for %s: %s", style.Render(m.filteredServicesList[m.pointer].name), m.snapshotDataTextBox.View())
	return s
}
