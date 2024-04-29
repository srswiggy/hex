package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Service struct {
	name     string
	selected bool
	input    string
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
