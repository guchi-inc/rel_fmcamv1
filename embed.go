// static/embed.go
package fmcam

import "embed"

//go:embed assets/* icons/* common/configs/* *.html *.svg *.json *.yaml
var FS embed.FS
