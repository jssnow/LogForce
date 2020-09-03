package common

import (
	"bytes"
	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

// output banner
func InitBanner() {
	bannerStr := `
______               __________
___  / _____________ ___  ____/_______________________
__  /  _  __ \_  __ /_  /_   _  __ \_  ___/  ___/  _ \
_  /___/ /_/ /  /_/ /_  __/   / /_/ /  /   / /__ /  __/
/_____/\____/_\__, / /_/      \____//_/    \___/ \___/
             /____/

GoVersion: {{ .GoVersion }}
GOOS: {{ .GOOS }}
GOARCH: {{ .GOARCH }}
NumCPU: {{ .NumCPU }}
GOPATH: {{ .GOPATH }}
GOROOT: {{ .GOROOT }}
Compiler: {{ .Compiler }}
ENV: {{ .Env "GOPATH" }}
Now: {{ .Now "Monday, 2 Jan 2006" }}
`
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString(bannerStr))
}
