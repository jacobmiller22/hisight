package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jacobmiller22/hisight/internal/config"
	"github.com/jacobmiller22/hisight/internal/logkeys"

	"github.com/jacobmiller22/gossentials/clog"
)

const hsHookUsage string = "Usage:\n\thshook [shell] [args]"

var ErrUnsupportedHook error = errors.New("unsupported hook")

type HookContext struct {
	SelfPath string
}

func HsHook(ctx context.Context, args []string) error {

	l := clog.FromContext(ctx)
	cfg := config.LoadConfig(args)

	l.Debug(logkeys.CommandStart, logkeys.Command, "HSHOOK", logkeys.Config, cfg)

	if len(args) < 1 {
		return fmt.Errorf(hsHookUsage)
	}
	target := args[0]

	sh := DetectShell(target)

	if sh == nil {
		return fmt.Errorf("%w: %s not supported", ErrUnsupportedHook, target)
	}

	hookStr, err := sh.Hook(strings.Join(args[1:], " "))
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
	Name() string
	Hook(args string) (string, error)
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
