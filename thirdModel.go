package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func (m thirdModel) Init() tea.Cmd {
	return nil
}

type thirdModel struct {
	servicesList []*Service
	envName      string
	jsonCreated  bool
	deployed     bool
	output       string
}

func (m thirdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if m.jsonCreated {
				m.deployed = true
				m.output = deployJson()
				//return m, tea.Quit
			}
			createAndSaveJson(m.envName, m.servicesList)
			m.jsonCreated = true
		} else if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else {
			m.envName += msg.String()
		}
	}

	return m, nil
}

func (m thirdModel) View() string {
	s := ""
	if m.deployed {
		return m.output
	} else if m.jsonCreated {
		s += fmt.Sprintf("\nPress Enter to trigger deployment\n")
	} else {
		s += fmt.Sprintf("Enter Env Name and press enter to create deploy.json\n Env Name: %s", m.envName)
	}
	return s
}
