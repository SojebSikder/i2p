package converter

// Postman structures

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
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol,omitempty"`
	Port     string   `json:"port,omitempty"`
	Host     []string `json:"host,omitempty"`
	Path     []string `json:"path,omitempty"`
}
