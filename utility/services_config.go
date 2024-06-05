package utility

import (
	"encoding/json"
	"fmt"
	"hex/data_model"
	"io/ioutil"
	"net/http"
)

func GetTemplates() data_model.Template {
	//cfg := config.GetConfig()
	url := "https://raw.githubusercontent.com/srswiggy/shuttel-services/main/services.json"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch JSON: %v\n", err)
		fmt.Printf(url)
		return data_model.Template{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch JSON: HTTP %d\n", resp.StatusCode)
		return data_model.Template{}
	}

	var templates data_model.Templates
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return data_model.Template{}
	}

	err = json.Unmarshal(body, &templates)
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return data_model.Template{}
	}

	return templates.Templates[0]
}
