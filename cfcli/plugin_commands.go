package cfcli

import (
	"fmt"
	"os"
	"strings"
)

var cfShellCmds = [][]string{
	{"quit", "cf-shell: Exits the shell"},
	{"exit", "cf-shell: Exits the shell"},
	// {"ls", "cf-shell: ls ..."},
	// {"dir", "cf-shell: dir ..‚"},
	// {"pwd", "cf-shell: pwd"},
	// {"cd", "cf-shell: cd"},
	// {"cls", "cf-shell: cls"},
}

func isPluginCommand(s string) bool {
	parts := strings.Split(s, " ")
	for _, cmd := range cfShellCmds {
		if parts[0] == cmd[0] {
			return true
		}
	}
	return false
}

func executePluginCommand(s string) {
	parts := strings.Split(s, " ")
	switch parts[0] {
	case "quit", "exit":
		fmt.Println("exiting cf-shell")
		os.Exit(0)
	default:
		executeShellCommand(s)
	}
}
