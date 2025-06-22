package main

import (
	"os"
	"text/template"
)

type HookContext struct {
	SelfPath string
}

func main() {

	target := "bash"

	sh := DetectShell(target)

	hookStr, err := sh.Hook()
	if err != nil {
		panic("error calling Hook()")
	}

	selfPath, err := os.Executable()
	if err != nil {
		panic("error calling os.Executable()")
	}

	hookContext := HookContext{
		SelfPath: selfPath,
	}

	hookTemplate, err := template.New("hook").Parse(hookStr)
	if err != nil {
		panic("error calling template")
	}

	hookTemplate.Execute(os.Stdout, hookContext)

}
