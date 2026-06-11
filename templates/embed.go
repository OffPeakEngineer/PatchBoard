package templates

import "embed"

// Files contains the project templates installed or rendered by Patchboard.
//
//go:embed *.html *.tmpl
var Files embed.FS
