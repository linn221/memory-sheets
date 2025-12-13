package views

import (
	"bytes"

	"github.com/yuin/goldmark"
)

const dateFormat = "2006-01-02"

// MarkdownToHTML converts markdown text to HTML
func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.New().Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MarkdownToHTMLSafe converts markdown text to HTML, returning empty string on error
func MarkdownToHTMLSafe(markdown string) string {
	html, err := MarkdownToHTML(markdown)
	if err != nil {
		return ""
	}
	return html
}
