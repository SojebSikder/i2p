package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var version = "0.0.1"

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  i2p convert [--input-file FILE] [--output-file FILE]")
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

// Input (Insomnia) structures

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
	Name     string             `yaml:"name"`
	URL      string             `yaml:"url,omitempty"`
	Method   string             `yaml:"method,omitempty"`
	Body     *RequestBody       `yaml:"body,omitempty"`
	Headers  []Header           `yaml:"headers,omitempty"`
	Children []InsomniaSubgroup `yaml:"children,omitempty"`
}

type RequestBody struct {
	MimeType string `yaml:"mimeType"`
	Text     string `yaml:"text"`
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `json:"value"`
}

// Output (Postman) structures

type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanInfo struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type PostmanItem struct {
	Name    string          `json:"name"`
	Request *PostmanRequest `json:"request,omitempty"`
	Item    []PostmanItem   `json:"item,omitempty"`
}

type PostmanRequest struct {
	Method string              `json:"method"`
	Header []Header            `json:"header,omitempty"`
	Body   *PostmanRequestBody `json:"body,omitempty"`
	URL    PostmanURL          `json:"url"`
}

type PostmanRequestBody struct {
	Mode    string                    `json:"mode"`
	Raw     string                    `json:"raw"`
	Options PostmanRequestBodyOptions `json:"options"`
}

type PostmanRequestBodyOptions struct {
	Raw PostmanRequestBodyOptionsRaw `json:"raw"`
}

type PostmanRequestBodyOptionsRaw struct {
	Language string `json:"language"`
}

type PostmanURL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host,omitempty"`
	Path []string `json:"path,omitempty"`
}

// Core conversion logic

func convertInsomniaToPostman(inputFile string, outputFile string) {
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
	path := strings.Split(strings.Trim(parsed.Path, "/"), "/")

	return PostmanURL{
		Raw:  rawURL,
		Host: host,
		Path: path,
	}
}
