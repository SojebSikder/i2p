package util

import (
	"net/url"
	"path"
)

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
