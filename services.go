package main

import (
	"encoding/json"
	"fmt"
	"hex/data_model"
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
	EnvName  string               `json:"envName"`
	Services []data_model.Service `json:"services"`
}

func createAndSaveJson(envName string, services []*struct {
	Service    data_model.Service
	IsSelected bool
	Input      string
}) {
	configuredServices := []data_model.Service{}
	for _, service := range services {
		service.Service.Version = service.Input
		configuredServices = append(configuredServices, service.Service)
	}

	newConfig := ConfigFile{
		EnvName:  envName,
		Services: configuredServices,
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

func deployJson(envname string) string {
	filename := "output.json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File does not exist: %s\n", filename)
		return "File does not exist"
	}

	// Create the qgp command with the filename
	cmd := exec.Command("qgp", "update", "instance", "--updateServices=true", "-n="+envname, "-f="+filename)

	// Execute the command and get the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputMessage := fmt.Sprintf("Failed to execute command: %s\n", err, string(output))
		return outputMessage
	}

	// Return the output
	return string(output)
}

func getServicesConfig(serviceName string) Services {
	var servicesConfigMap = make(map[string]Services)

	servicesConfigMap["finance-calcy-service"] = Services{
		Name:              "finance-calcy-service",
		DependentServices: make([]string, 0),
		IsMockService:     false,
		CommitSha:         "",
		File:              "app.yaml",
		BranchName:        "master",
		Repo:              "finance-calcy-service",
	}

	servicesConfigMap["finance-job-service"] = Services{
		Name:              "finance-job-service",
		DependentServices: []string{"finance-calcy-service"},
		IsMockService:     false,
		CommitSha:         "",
		File:              "app.yaml",
		BranchName:        "master",
		Repo:              "finance-job",
	}

	servicesConfigMap["finance-dashboard"] = Services{
		Name:              "finance-dashboard",
		DependentServices: []string{"finance-calcy-service", "lassi"},
		IsMockService:     false,
		CommitSha:         "",
		File:              "app.yaml",
		BranchName:        "master",
		Repo:              "finance-dashboard",
	}

	servicesConfigMap["finance-orchestrator"] = Services{
		Name:              "finance-orchestrator",
		DependentServices: []string{"finance-recon-platform", "finance-scheduler-service", "finance-calcy-service"},
		IsMockService:     false,
		CommitSha:         "",
		File:              "app.yaml",
		BranchName:        "master",
		Repo:              "finance-orchestrator",
	}

	return servicesConfigMap[serviceName]
}
