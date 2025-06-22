package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const usage string = "Usage:\n\thshook [shell] [args]"

var ErrUnsupportedHook error = errors.New("unsupported hook")

type HookContext struct {
	SelfPath string
}

func HsHook(args []string) error {

	if len(args) < 1 {
		return fmt.Errorf(usage)
	}
	target := args[0]

	sh := DetectShell(target)

	if sh == nil {
		return fmt.Errorf("%w: %s not supported", ErrUnsupportedHook, target)
	}

	hookStr, err := sh.Hook()
	if err != nil {
		return fmt.Errorf("error calling Hook()")
	}

	selfPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error calling os.Executable()")
	}

	hookContext := HookContext{
		SelfPath: selfPath,
	}

	hookTemplate, err := template.New("hook").Parse(hookStr)
	if err != nil {
		return fmt.Errorf("error calling template")
	}

	hookTemplate.Execute(os.Stdout, hookContext)
	return nil

}

type Shell interface {
	Hook() (string, error)
}

var supportedShellList = map[string]Shell{
	"bash": Bash,
	// "zsh":     Zsh,
}

// DetectShell returns a Shell instance from the given target.
//
// target is usually $0 and can also be prefixed by `-`
func DetectShell(target string) Shell {
	target = filepath.Base(target)
	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	detectedShell, isValid := supportedShellList[target]
	if isValid {
		return detectedShell
	}
	return nil
}
