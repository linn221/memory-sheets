package views

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const dateFormat = "2006-01-02"

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

// MarkdownToHTML converts markdown text to HTML
func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
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
