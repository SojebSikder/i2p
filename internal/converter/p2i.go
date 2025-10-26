package converter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ConvertPostmanToInsomnia(inputFile string, outputFile string) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var postman PostmanCollection
	if err := json.Unmarshal(data, &postman); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	insomnia := InsomniaExport{
		Type:       "export",
		Name:       postman.Info.Name,
		Collection: []InsomniaFolder{},
	}

	for _, item := range postman.Item {
		folder := convertPostmanItemToFolder(item)
		insomnia.Collection = append(insomnia.Collection, folder)
	}

	yamlData, err := yaml.Marshal(insomnia)
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
		log.Fatalf("Failed to write YAML file: %v", err)
	}

	fmt.Println("Conversion successful! Output:", outputFile)
}

func convertPostmanItemToFolder(item PostmanItem) InsomniaFolder {
	folder := InsomniaFolder{
		Name:     item.Name,
		Children: []InsomniaSubgroup{},
	}

	for _, child := range item.Item {
		if child.Request != nil {
			sub := InsomniaSubgroup{
				Name:    child.Name,
				URL:     reconstructURL(child.Request.URL),
				Method:  child.Request.Method,
				Headers: child.Request.Header,
			}

			if child.Request.Body != nil {
				sub.Body = &RequestBody{
					MimeType: "application/json",
					Text:     child.Request.Body.Raw,
				}
			}

			folder.Children = append(folder.Children, sub)
		} else if len(child.Item) > 0 {
			// Nested folder
			subFolder := convertPostmanItemToFolder(child)
			sub := InsomniaSubgroup{
				Name:     subFolder.Name,
				Children: subFolder.Children,
			}
			folder.Children = append(folder.Children, sub)
		}
	}

	return folder
}

func reconstructURL(u PostmanURL) string {
	if u.Raw != "" {
		return u.Raw
	}

	// fallback: reconstruct URL from parts
	urlStr := ""
	if u.Protocol != "" {
		urlStr += u.Protocol + "://"
	}

	if len(u.Host) > 0 {
		urlStr += strings.Join(u.Host, ".")
	}

	if u.Port != "" {
		urlStr += ":" + u.Port
	}

	if len(u.Path) > 0 {
		urlStr += "/" + strings.Join(u.Path, "/")
	}

	return urlStr
}
