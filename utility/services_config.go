package utility

import (
	"encoding/json"
	"fmt"
	"hex/data_model"
	"io/ioutil"
	"net/http"
)

func GetTemplates() data_model.Templates {
	//cfg := config.GetConfig()
	url := "https://raw.githubusercontent.com/srswiggy/shuttel-services/main/services.json"
	config, err := data_model.LoadConfig()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch JSON: %v\n", err)
		fmt.Printf(url)
		return data_model.Templates{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch JSON: HTTP %d\n", resp.StatusCode)
		return data_model.Templates{}
	}

	var templates data_model.Templates
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return data_model.Templates{}
	}

	err = json.Unmarshal(body, &templates)
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return data_model.Templates{}
	}
	return data_model.Templates{
		Templates: Filter(templates.Templates, func(t data_model.Template) bool {
			if len(config.GetSelectedPods()) == 0 {
				return true
			}
			for _, pod := range config.GetSelectedPods() {
				if t.Pod == pod {
					return true
				}
			}
			return false
		}),
	}
}
