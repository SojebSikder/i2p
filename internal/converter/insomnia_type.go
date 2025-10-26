package converter

// Insomnia structures

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
