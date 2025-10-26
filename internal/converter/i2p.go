package converter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ConvertInsomniaToPostman(inputFile string, outputFile string) {
	yamlData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var insomnia InsomniaExport
	if err := yaml.Unmarshal(yamlData, &insomnia); err != nil {
		log.Fatalf("YAML unmarshal failed: %v", err)
	}

	postman := PostmanCollection{
		Info: PostmanInfo{
			Name:   insomnia.Name,
			Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
	}

	for _, folder := range insomnia.Collection {
		item := convertGroup(folder.Name, folder.Children)
		postman.Item = append(postman.Item, item)
	}

	output, err := json.MarshalIndent(postman, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal failed: %v", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		log.Fatalf("Write file failed: %v", err)
	}

	fmt.Println("Conversion successful! Output:", outputFile)
}

func convertGroup(name string, children []InsomniaSubgroup) PostmanItem {
	item := PostmanItem{Name: name}

	for _, child := range children {
		if child.URL != "" && child.Method != "" {
			// It's a request
			var body *PostmanRequestBody
			if child.Body != nil {
				body = &PostmanRequestBody{
					Mode: "raw",
					Raw:  child.Body.Text,
					Options: PostmanRequestBodyOptions{
						Raw: PostmanRequestBodyOptionsRaw{
							Language: "json",
						},
					},
				}
			}

			req := PostmanItem{
				Name: child.Name,
				Request: &PostmanRequest{
					Method: child.Method,
					Header: child.Headers,
					Body:   body,
					URL:    parseURL(child.URL),
				},
			}
			item.Item = append(item.Item, req)
		} else {
			// It's a folder
			subItem := convertGroup(child.Name, child.Children)
			item.Item = append(item.Item, subItem)
		}
	}

	return item
}

func cleanTemplateVar(input string) string {
	input = strings.ReplaceAll(input, "_['", "")
	input = strings.ReplaceAll(input, "']", "")
	return input
}

func parseURL(rawURL string) PostmanURL {
	if strings.HasPrefix(rawURL, "{{") {
		// Handle variable-style URL like {{endpoint}}/api/auth/register
		parts := strings.SplitN(rawURL, "/", 2)
		host := []string{cleanTemplateVar(parts[0])}
		path := []string{}
		if len(parts) > 1 {
			path = strings.Split(parts[1], "/")
		}
		return PostmanURL{
			Raw:  rawURL,
			Host: host,
			Path: path,
		}
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return PostmanURL{Raw: rawURL}
	}

	host := strings.Split(parsed.Hostname(), ".")
	// remove protocol from host
	host = strings.Split(host[0], ":")
	host = host[:len(host)-1]
	host = append(host, parsed.Hostname())

	port := parsed.Port()

	path := strings.Split(strings.Trim(parsed.Path, "/"), "/")

	return PostmanURL{
		Raw:      rawURL,
		Protocol: parsed.Scheme,
		Host:     host,
		Port:     port,
		Path:     path,
	}
}
