package main

import "path/filepath"

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
