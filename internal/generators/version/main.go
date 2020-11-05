package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/consensys/bavard"
)

//go:generate go run main.go
func main() {
	v, err := exec.Command("git", "describe", "--abbrev=0").CombinedOutput()
	if err != nil {
		panic(err)
	}
	version := strings.TrimSpace(string(v))
	src := []string{
		versionTemplate,
	}

	if err := bavard.Generate("../../../generator/version.go", src,
		struct{ Version string }{version},
		bavard.Apache2("ConsenSys Software Inc.", 2020),
		bavard.Package("generator"),
		bavard.GeneratedBy("internal/generators/version")); err != nil {
		fmt.Println("error", err)
		os.Exit(-1)
	}
}

const versionTemplate = `
// Version goff version
const Version = "{{.Version}}"
`
