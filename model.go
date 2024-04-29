package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type model struct {
	servicesList           []*Service
	pointer                int
	filteredServicesList   []*Service
	listSearchQueryTextBox textinput.Model
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 20
	ti.Placeholder = "Search Service"

	servicesList := []*Service{
		{name: "finance-calcy-service", selected: false},
		{name: "finance-job-service", selected: false},
		{name: "finance-orchestrator", selected: false},
		{name: "finance-dashboard", selected: false},
	}
	return model{
		servicesList:           servicesList,
		filteredServicesList:   []*Service{},
		listSearchQueryTextBox: ti,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if msg.Type == tea.KeyUp {
			if m.pointer > 0 {
				m.pointer--
			}
		} else if msg.Type == tea.KeyDown {
			if m.pointer < len(m.filteredServicesList)-1 {
				m.pointer++
			}
		} else if msg.Type == tea.KeyEnter {
			var selectedServices []*Service
			for _, service := range m.servicesList {
				if service.selected {
					selectedServices = append(selectedServices, service)
				}
			}
			if len(selectedServices) == 0 {
				return m, cmd
			}
			return secondModel{filteredServicesList: selectedServices}, nil
		} else if msg.String() == " " || msg.Type == tea.KeySpace {
			if m.filteredServicesList[m.pointer].selected {
				m.filteredServicesList[m.pointer].selected = false
			} else {
				m.filteredServicesList[m.pointer].selected = true
			}
		} else {
			// Update the listSearchQuery with the new character
			m.listSearchQueryTextBox, cmd = m.listSearchQueryTextBox.Update(msg)
			//m.listSearchQuery += msg.String()
		}

		// Filter the items based on the listSearchQuery
		m.filteredServicesList = []*Service{}
		if len(m.listSearchQueryTextBox.Value()) > 0 {
			for _, item := range m.servicesList {
				if strings.Contains(strings.ToLower(item.name), strings.ToLower(m.listSearchQueryTextBox.Value())) {
					m.filteredServicesList = append(m.filteredServicesList, item)
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("5"))

	s := "Enter your search query:\n"
	s += m.listSearchQueryTextBox.View() + "\n\n"

	if len(m.filteredServicesList) == 0 {
		s += "No matching items found"
	} else {
		for i, item := range m.filteredServicesList {
			pointer := " "
			if m.pointer == i {
				pointer = ">"
			}

			checked := " "
			if item.selected {
				checked = "x"
			}

			if pointer == ">" {
				s += style.Render(fmt.Sprintf("%s [%s] %s", pointer, checked, item.name)) + "\n"
			} else {
				s += fmt.Sprintf("%s [%s] %s", pointer, checked, item.name) + "\n"
			}

		}
	}
	// Show Selected Services
	selectedServices := "\n\n============ Selected Services ==========\n"
	i := 1
	for _, service := range m.servicesList {
		if service.selected == true {
			selectedServices += fmt.Sprintf("%d. %s\n", i, service.name)
			i++
		}
	}
	s += selectedServices
	return s
}
