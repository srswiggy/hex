package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type secondModel struct {
	filteredServicesList []*Service
	pointer              int
	snapshotData         string
}

func (m secondModel) Init() tea.Cmd {
	return nil
}

func (m secondModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else if msg.Type == tea.KeyEnter {
			if m.pointer >= len(m.filteredServicesList) {
				return thirdModel{servicesList: m.filteredServicesList}, nil
			}
			m.filteredServicesList[m.pointer].input = m.snapshotData
			m.snapshotData = ""
			m.pointer++
		} else {
			m.snapshotData += msg.String()
		}
	}
	return m, nil
}

func (m secondModel) View() string {
	s := ""
	if m.pointer >= len(m.filteredServicesList) {
		return "All inputs complete, Press Enter to create deploy.json"
	}
	for _, service := range m.filteredServicesList {
		if len(service.input) > 0 {
			s += fmt.Sprintf("%s: %s\n", service.name, service.input)
		}
	}
	s += fmt.Sprintf("Input data for %s: %s", m.filteredServicesList[m.pointer].name, m.snapshotData)
	return s
}
