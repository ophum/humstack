package apply

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

func printYAML(v interface{}) {
	yamlBytes, err := yaml.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(yamlBytes))
}
