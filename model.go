package main

import (
	"fmt"
	"hex/data_model"
	"hex/utility"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	servicesList           []*data_model.ModelService
	pointer                int
	listSearchQueryTextBox textinput.Model
	filteredServicesList   []*data_model.ModelService
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func initialModel() model {
	fetchedEnvironmentTemplates := utility.GetTemplates()
	services := initializeServicesList(&fetchedEnvironmentTemplates)
	// fmt.Println(services)
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
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
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
			var selectedServices []*data_model.ModelService
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
		m.filteredServicesList = []*data_model.ModelService{}
		if len(m.listSearchQueryTextBox.Value()) > 0 {
			for _, item := range m.servicesList {
				if strings.Contains(strings.ToLower(item.Service.Name), strings.ToLower(m.listSearchQueryTextBox.Value())) {
					m.filteredServicesList = append(m.filteredServicesList, item)
				}
			}
		} else if len(m.listSearchQueryTextBox.Value()) == 0 {
			m.filteredServicesList = m.servicesList
		}

		// Fix for Issue 2: Reset pointer if it's out of bounds after filtering
		if len(m.filteredServicesList) == 0 {
			m.pointer = 0
		} else if m.pointer >= len(m.filteredServicesList) {
			m.pointer = 0
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("5"))
	// var greenStyle = lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("3"))
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

		var mockStatus string = ""
		if item.Service.IsMockService {
			mockStatus = "[ Mock Service ]"
		}

		if pointer == ">" {
			s += style.Render(fmt.Sprintf("%s [%s] %-30s [ %s ] %s", pointer, checked, item.Service.Name, item.Pod, mockStatus))
			s += "\n"
		} else {
			s += fmt.Sprintf("%s [%s] %-30s [ %s ] %s\n", pointer, checked, item.Service.Name, item.Pod, mockStatus)
		}
	}

	if listlen == 0 {
		s += fmt.Sprintf("No results found in your search query.\n\n")
	}

	s += firstScreenBottomIntructions()
	return s
}

func initializeServicesList(templates *data_model.Templates) []*data_model.ModelService {
	var newList []*data_model.ModelService
	for _, template := range templates.Templates {
		services := template.Services

		for _, service := range services {
			newList = append(newList, &data_model.ModelService{
				Service:    service,
				IsSelected: false,
				Input:      "",
				Pod:        template.Pod,
			})
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

func firstScreenBottomIntructions() string {
	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("5"))
	return "\n\n" +
		style.Render("Space Bar") +
		" to select/deselect service\n" +
		style.Render("Ctrl+C or Esc") +
		" to quit anytime\n" +
		style.Render("Enter") +
		" to proceed to another screen"
}
