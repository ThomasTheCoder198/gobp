package render

import "embed"

// templatesFS embeds every file under internal/templates as the project
// template tree. The leading "templates/" prefix is preserved in paths.
//
//go:embed all:templates
var templatesFS embed.FS
