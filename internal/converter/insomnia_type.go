package converter

// Insomnia structures

type InsomniaExport struct {
	Type        string                 `yaml:"type"`
	Name        string                 `yaml:"name"`
	Meta        Meta                   `yaml:"meta"`
	Collection  []InsomniaFolder       `yaml:"collection"`
	Environment map[string]interface{} `yaml:"environments"`
}

type Meta struct {
	ID       string `yaml:"id"`
	Created  int64  `yaml:"created"`
	Modified int64  `yaml:"modified"`
}

type FolderMeta struct {
	ID       string `yaml:"id"`
	Created  int64  `yaml:"created"`
	Modified int64  `yaml:"modified"`
	SortKey  int64  `yaml:"sortKey"`
}

type InsomniaFolder struct {
	Name     string             `yaml:"name"`
	Meta     FolderMeta         `yaml:"meta"`
	Children []InsomniaSubgroup `yaml:"children"`
}

type InsomniaSubgroup struct {
	URL      string             `yaml:"url,omitempty"`
	Name     string             `yaml:"name"`
	Meta     FolderMeta         `yaml:"meta"`
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
