package util

import (
	"bytes"
	"net/url"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

func MarshalYAML(v any) ([]byte, error) {
	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	defer enc.Close()

	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetTime() int64 {
	now := time.Now()
	return now.Unix()
}

func GetFileNameFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "file"
	}
	if name := u.Query().Get("filename"); name != "" {
		return name
	}
	return path.Base(u.Path)
}

func GetFileExtensionFromURL(rawURL string) string {
	fileName := GetFileNameFromURL(rawURL)
	return path.Ext(fileName)
}
