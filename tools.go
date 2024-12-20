//go:build tools
// +build tools

package tools

// Go does not have native tool management
// Therefore we need to hack a little bit
// https://play-with-go.dev/tools-as-dependencies_go119_en/
import (
	_ "github.com/a-h/templ/cmd/templ"
)
