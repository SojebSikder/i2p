package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type InsomniaExport struct {
	Type        string                 `yaml:"type"`
	Name        string                 `yaml:"name"`
	Collection  []InsomniaFolder       `yaml:"collection"`
	Environment map[string]interface{} `yaml:"environment"`
}

type InsomniaFolder struct {
	Name     string             `yaml:"name"`
	Children []InsomniaSubgroup `yaml:"children"`
}

type InsomniaSubgroup struct {
	Name     string            `yaml:"name"`
	Children []InsomniaRequest `yaml:"children"`
}

type InsomniaRequest struct {
	Name    string       `yaml:"name"`
	URL     string       `yaml:"url"`
	Method  string       `yaml:"method"`
	Body    *RequestBody `yaml:"body"`
	Headers []Header     `yaml:"headers"`
}

type RequestBody struct {
	MimeType string `yaml:"mimeType"`
	Text     string `yaml:"text"`
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Postman v2.1 structures

type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanInfo struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type PostmanItem struct {
	Name    string         `json:"name"`
	Request PostmanRequest `json:"request"`
}

type PostmanRequest struct {
	Method string              `json:"method"`
	Header []Header            `json:"header"`
	Body   *PostmanRequestBody `json:"body,omitempty"`
	URL    PostmanURL          `json:"url"`
}

type PostmanRequestBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type PostmanURL struct {
	Raw string `json:"raw"`
}

var version = "0.0.1"

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  i2p convert [--input-file FILE] [--output-file FILE]")
	fmt.Println()
	fmt.Println("  i2p help")
	fmt.Println("  i2p version")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --input-file FILE   Specify the input file (default: insomnia.yaml)")
	fmt.Println("  --output-file FILE  Specify the output file (default: postman_collection.json)")
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "convert":
		inputFile := "insomnia.yaml"
		outputFile := "postman_collection.json"

		fs := flag.NewFlagSet("convert", flag.ExitOnError)
		fs.StringVar(&inputFile, "input-file", inputFile, "Specify the input file")
		fs.StringVar(&outputFile, "output-file", outputFile, "Specify the output file")
		fs.Parse(os.Args[2:])

		convertInsomniaToPostman(inputFile, outputFile)

	case "help":
		showUsage()
	case "version":
		fmt.Println("i2p version " + version)
	default:
		fmt.Println("Unknown command:", cmd)
		fmt.Println("Use 'i2p help' to see available commands.")
		os.Exit(1)
	}
}

func convertInsomniaToPostman(inputFile string, outputFile string) {
	// Load YAML file
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

	// Convert Insomnia to Postman
	for _, folder := range insomnia.Collection {
		for _, group := range folder.Children {
			for _, req := range group.Children {
				var body *PostmanRequestBody
				if req.Body != nil {
					body = &PostmanRequestBody{
						Mode: "raw",
						Raw:  req.Body.Text,
					}
				}

				item := PostmanItem{
					Name: req.Name,
					Request: PostmanRequest{
						Method: req.Method,
						Header: req.Headers,
						Body:   body,
						URL:    PostmanURL{Raw: req.URL},
					},
				}
				postman.Item = append(postman.Item, item)
			}
		}
	}

	// Write to Postman JSON
	output, err := json.MarshalIndent(postman, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal failed: %v", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		log.Fatalf("Write file failed: %v", err)
	}

	fmt.Println("Conversion successful! Output: " + outputFile)
}
