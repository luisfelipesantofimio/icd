package model

import (
	"encoding/json"
	"fmt"
	"os"
)

// Json serializable instance of an ICD entity
type IcdBlock struct {
	Title struct {
		Value string `json:"@value"`
	} `json:"title"`
	Definition struct {
		Value string `json:"@value"`
	} `json:"definition,omitempty"`
	BlockType  string `json:"classKind"`
	BrowserURL string `json:"browserUrl,omitempty"`
	Code       string `json:"code,omitempty"`
	Children   []IcdBlock
}

// TODO: Create repo
func ExportToJSON(model *IcdBlock, locale string) error {
	data, _ := json.Marshal(model)

	outputDir := "output"

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	file := "./" + outputDir + "/" + "icd_data_" + locale + ".json"
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
