package main

import (
	"github.com/russross/blackfriday"
)

const (
	flags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_SMARTYPANTS_ANGLED_QUOTES
	extensions = 0 |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_TABLES
)

func renderMarkdown(md []byte) (html string) {
	renderer := blackfriday.HtmlRenderer(flags, "", "")
	html = string(blackfriday.Markdown(md, renderer, extensions))
	return
}
