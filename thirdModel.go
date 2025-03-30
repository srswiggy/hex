package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"hex/data_model"
	"hex/utility"
)

func (m thirdModel) Init() tea.Cmd {
	return nil
}

type thirdModel struct {
	servicesList   []*data_model.ModelService
	envNameTextBox textinput.Model
	jsonCreated    bool
	deployed       bool
	output         string
}

func (m thirdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if m.jsonCreated {
				m.deployed = true
				m.output = deployJson(m.envNameTextBox.Value())
			}
			createAndSaveJson(m.envNameTextBox.Value(), m.servicesList)
			m.jsonCreated = true
		} else if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else {
			m.envNameTextBox, cmd = m.envNameTextBox.Update(msg)
		}
	}

	return m, cmd
}

func (m thirdModel) View() string {
	s := ""
	if m.deployed {
		return wordwrap.String(m.output, utility.GetTerminalWidth())
	} else if m.jsonCreated {
		s += fmt.Sprintf("\nPress Enter to trigger deployment\n")
	} else {
		s += fmt.Sprintf("Enter Env Name and press enter to create deploy.json\n Env Name: %s", m.envNameTextBox.View())
	}
	return s
}
