package parser

import (
	"gopkg.in/yaml.v3"
	"oncall-gen/internal/model"
	"os"
)

func ParseYaml(filename string) (*model.Teams, error) {
	var result model.Teams

	file, err := os.ReadFile(filename)
	if err != nil {
		return &result, err
		//log.Fatalf("Error while reading input file %s: %v", filename, err)
	}

	err = yaml.Unmarshal(file, &result)
	if err != nil {
		return &result, err
		//log.Fatalf("Error while parsing file content: %v", err)
	}

	return &result, nil
}
