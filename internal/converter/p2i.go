package converter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sojebsikder/i2p/internal/util"
	"strings"

	"github.com/google/uuid"
)

var created = util.GetTime()
var modified = util.GetTime()

const sortKey = -1754819949124

func ConvertPostmanToInsomnia(inputFile string, outputFile string) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var postman PostmanCollection
	if err := json.Unmarshal(data, &postman); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	id := uuid.New()

	insomnia := InsomniaExport{
		Type: "collection.insomnia.rest/5.0",
		Name: postman.Info.Name,
		Meta: Meta{
			ID:       id.String(),
			Created:  created,
			Modified: modified,
		},
		Environment: ConvertVariablesToEnv(postman.Variable),
		Collection:  []InsomniaFolder{},
	}

	for _, item := range postman.Item {
		folder := convertPostmanItemToFolder(item)
		insomnia.Collection = append(insomnia.Collection, folder)
	}

	yamlData, err := util.MarshalYAML(insomnia)

	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
		log.Fatalf("Failed to write YAML file: %v", err)
	}

	fmt.Println("Conversion successful! Output:", outputFile)
}

func convertPostmanItemToFolder(item PostmanItem) InsomniaFolder {
	id := uuid.New()
	folder := InsomniaFolder{
		Name: item.Name,
		Meta: FolderMeta{
			ID:       id.String(),
			Created:  created,
			Modified: modified,
			SortKey:  sortKey,
		},
		Children: []InsomniaSubgroup{},
	}

	for _, child := range item.Item {
		if child.Request != nil {
			id := uuid.New()
			sub := InsomniaSubgroup{
				Name: child.Name,
				URL:  reconstructURL(child.Request.URL),
				Meta: FolderMeta{
					ID:       id.String(),
					Created:  created,
					Modified: modified,
					SortKey:  sortKey,
				},
				Method:  child.Request.Method,
				Headers: child.Request.Header,
			}

			if child.Request.Body != nil {
				bodyText := strings.TrimSpace(child.Request.Body.Raw)

				var prettyJSON map[string]interface{}
				if json.Unmarshal([]byte(bodyText), &prettyJSON) == nil {
					formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
					bodyText = string(formatted)
				}

				sub.Body = &RequestBody{
					MimeType: "application/json",
					Text:     bodyText,
				}
			}

			folder.Children = append(folder.Children, sub)
		} else if len(child.Item) > 0 {
			subFolder := convertPostmanItemToFolder(child)
			id := uuid.New()
			sub := InsomniaSubgroup{
				Name: subFolder.Name,
				Meta: FolderMeta{
					ID:       id.String(),
					Created:  created,
					Modified: modified,
					SortKey:  sortKey,
				},
				Children: subFolder.Children,
			}
			folder.Children = append(folder.Children, sub)
		}
	}

	return folder
}

func ConvertVariablesToEnv(vars []Variable) map[string]interface{} {
	env := map[string]interface{}{}
	data := map[string]string{}

	for _, v := range vars {
		if v.Key == "endpoint" {
			data["base_url"] = v.Value
		} else {
			data[v.Key] = v.Value
		}
	}

	env["data"] = data
	return env
}

func reconstructURL(u PostmanURL) string {
	if u.Raw != "" {
		// Handle Postman variable templates {{variable}}
		if strings.HasPrefix(u.Raw, "{{") && strings.Contains(u.Raw, "}}") {
			parts := strings.SplitN(u.Raw, "/", 2)
			host := cleanTemplateVarP2i(parts[0])
			path := []string{}
			if len(parts) > 1 {
				path = strings.Split(parts[1], "/")
			}
			return host + "/" + strings.Join(path, "/")
		}
		return u.Raw
	}

	// Fallback: reconstruct URL from parts
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

func cleanTemplateVarP2i(input string) string {
	// Postman variable {{variable}} -> Insomnia {{variable}}
	input = strings.TrimSpace(input)
	input = strings.TrimPrefix(input, "{{")
	input = strings.TrimSuffix(input, "}}")
	return "{{ _." + input + " }}"
}
