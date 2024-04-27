package main

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"strings"
)

type Service struct {
	name     string
	selected bool
	input    string
}
type model struct {
	servicesList         []*Service
	pointer              int
	filteredServicesList []*Service
	listSearchQuery      string
}

func waitForInput(prompt string) tea.Cmd {
	return func() tea.Msg {
		fmt.Println(prompt + ": ")
		var input string
		fmt.Scanln(&input)
		return input
	}
}

func initialModel() model {
	servicesList := []*Service{
		{name: "finance-calcy-service", selected: false},
		{name: "finance-job-service", selected: false},
		{name: "finance-orchestrator", selected: false},
		{name: "finance-dashboard", selected: false},
	}
	return model{
		servicesList:         servicesList,
		filteredServicesList: []*Service{},
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m secondModel) Init() tea.Cmd {
	return nil
}

func (m thirdModel) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// Handle backspace
		if msg.Type == tea.KeyBackspace {
			if len(m.listSearchQuery) > 0 {
				m.listSearchQuery = m.listSearchQuery[:len(m.listSearchQuery)-1]
			}
		} else if msg.Type == tea.KeyUp {
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
			return secondModel{filteredServicesList: selectedServices}, nil
		} else if msg.String() == " " {
			if m.filteredServicesList[m.pointer].selected {
				m.filteredServicesList[m.pointer].selected = false
			} else {
				m.filteredServicesList[m.pointer].selected = true
			}
		} else {
			// Update the listSearchQuery with the new character
			m.listSearchQuery += msg.String()
		}

		// Filter the items based on the listSearchQuery
		m.filteredServicesList = []*Service{}
		if len(m.listSearchQuery) > 0 {
			for _, item := range m.servicesList {
				if strings.Contains(strings.ToLower(item.name), strings.ToLower(m.listSearchQuery)) {
					m.filteredServicesList = append(m.filteredServicesList, item)
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

type secondModel struct {
	filteredServicesList []*Service
	pointer              int
	snapshotData         string
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

func (m model) View() string {
	s := "Enter your search query:\n"
	s += m.listSearchQuery + "\n\n"

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

			s += fmt.Sprintf("%s [%s] %s\n", pointer, checked, item.name)
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

type Services struct {
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	DependentServices []string `json:"dependentServices"`
	IsMockService     bool     `json:"isMockService"`
	CommitSha         string   `json:"commitSha"`
	File              string   `json:"file"`
	BranchName        string   `json:"branchName"`
	Repo              string   `json:"repo"`
}
type ConfigFile struct {
	EnvName  string     `json:"envName"`
	Services []Services `json:"services"`
}

func createAndSaveJson(envName string, services []*Service) {
	var configServices []Services
	for _, service := range services {
		configServices = append(configServices, Services{
			Name:              service.name,
			Version:           service.input,
			DependentServices: make([]string, 0),
			IsMockService:     false,
			CommitSha:         "",
			File:              "app.yaml",
			BranchName:        "",
			Repo:              "",
		})
	}
	newConfig := &ConfigFile{
		EnvName:  envName,
		Services: configServices,
	}

	jsonData, err := json.Marshal(newConfig)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}
}

func deployJson() string {
	filename := "output.json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File does not exist: %s\n", filename)
		return "File does not exist"
	}

	// Create the cat command with the filename
	cmd := exec.Command("cat", filename)

	// Execute the command and get the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		return "Failed to execute command."
	}

	// Return the output
	return string(output)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
