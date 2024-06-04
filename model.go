package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hex/data_model"
	"hex/utility"
	"strings"
)

type model struct {
	servicesList []*struct {
		Service    data_model.Service
		IsSelected bool
		Input      string
	}
	pointer                int
	listSearchQueryTextBox textinput.Model
	filteredServicesList   []*struct {
		Service    data_model.Service
		IsSelected bool
		Input      string
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func initialModel() model {
	fetchedTemplates := utility.GetTemplates()
	services := initializeServicesList(fetchedTemplates.Services)
	return model{
		servicesList:           services,
		listSearchQueryTextBox: generateTextInputBox(),
		filteredServicesList:   services,
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
			var selectedServices []*struct {
				Service    data_model.Service
				IsSelected bool
				Input      string
			}
			for _, service := range m.servicesList {
				if service.IsSelected {
					selectedServices = append(selectedServices, service)
				}
			}
			if len(selectedServices) == 0 {
				return m, cmd
			}
			ti := textinput.New()
			ti.Focus()
			ti.Placeholder = "Enter Snapshot Tag for Service"
			return secondModel{filteredServicesList: selectedServices, snapshotDataTextBox: ti}, nil
		} else if msg.String() == " " || msg.Type == tea.KeySpace {
			if m.filteredServicesList[m.pointer].IsSelected {
				m.filteredServicesList[m.pointer].IsSelected = false
			} else {
				m.filteredServicesList[m.pointer].IsSelected = true
			}
		} else {
			// Update the listSearchQuery with the new character
			m.listSearchQueryTextBox, cmd = m.listSearchQueryTextBox.Update(msg)
			//m.listSearchQuery += msg.String()
		}

		// Filter the items based on the listSearchQuery
		m.filteredServicesList = []*struct {
			Service    data_model.Service
			IsSelected bool
			Input      string
		}{}
		if len(m.listSearchQueryTextBox.Value()) > 0 {
			for _, item := range m.servicesList {
				if strings.Contains(strings.ToLower(item.Service.Name), strings.ToLower(m.listSearchQueryTextBox.Value())) {
					m.filteredServicesList = append(m.filteredServicesList, item)
				}
			}
		} else if len(m.listSearchQueryTextBox.Value()) == 0 {
			m.filteredServicesList = m.servicesList
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

	//displayThreshold := config.GetConfig().DisplayThreshold
	displayThreshold := 6
	var displayStartIndex, displayEndIndex int

	// Calculate start and end indices
	listlen := len(m.filteredServicesList)
	if m.pointer < displayThreshold/2 {
		displayStartIndex = 0
		displayEndIndex = displayThreshold
	} else if m.pointer >= listlen-displayThreshold/2 {
		displayStartIndex = listlen - displayThreshold
		displayEndIndex = listlen
	} else {
		displayStartIndex = m.pointer - displayThreshold/2
		displayEndIndex = m.pointer + displayThreshold/2
	}

	// Adjust displayEndIndex if it's out of range
	if displayEndIndex > listlen {
		displayEndIndex = listlen
	}

	// Render the list
	for i := displayStartIndex; i < displayEndIndex; i++ {
		item := m.filteredServicesList[i]
		pointer := " "
		if m.pointer == i {
			pointer = ">"
		}

		checked := " "
		if i < listlen && m.filteredServicesList[i].IsSelected {
			checked = "x"
		}
		if pointer == ">" {
			s += style.Render(fmt.Sprintf("%s [%s] %s", pointer, checked, item.Service.Name))
			s += "\n"
		} else {
			s += fmt.Sprintf("%s [%s] %s\n", pointer, checked, item.Service.Name)
		}
	}

	if listlen == 0 {
		s += fmt.Sprintf("No results found in your search query.\n\n")
	}
	return s
}

func initializeServicesList(services []data_model.Service) []*struct {
	Service    data_model.Service
	IsSelected bool
	Input      string
} {
	newList := make([]*struct {
		Service    data_model.Service
		IsSelected bool
		Input      string
	}, len(services))

	for i := range newList {
		newList[i] = &struct {
			Service    data_model.Service
			IsSelected bool
			Input      string
		}{
			Service:    services[i],
			IsSelected: false,
			Input:      "",
		}
	}

	return newList
}

func generateTextInputBox() textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 20
	ti.Placeholder = "Search Service"
	return ti
}
