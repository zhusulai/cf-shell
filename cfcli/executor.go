package cfcli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var windowsInternals = []string{
	"ASSOC",
	"BREAK",
	"CALL",
	// "CD",
	"CHDIR",
	"CLS",
	"COLOR",
	"COPY",
	"DATE",
	"DEL",
	"DIR",
	"DPATH",
	"ECHO",
	"ENDLOCAL",
	"ERASE",
	"EXIT",
	"FOR",
	"FTYPE",
	"GOTO",
	"IF",
	"KEYS",
	"MD",
	"MKDIR",
	"MKLINK",
	"MOVE",
	"PATH",
	"PAUSE",
	"POPD",
	"PROMPT",
	"PUSHD",
	"REM",
	"REN",
	"RENAME",
	"RD",
	"RMDIR",
	"SET",
	"SETLOCAL",
	"SHIFT",
	"START",
	"TIME",
	"TITLE",
	"TYPE",
	"VER",
	"VERIFY",
	"VOL",
}

func isWindowsInternalCommand(s string) bool {
	parts := strings.Split(s, " ")
	for _, cmd := range windowsInternals {
		if strings.EqualFold(parts[0], cmd) {
			return true
		}
	}
	return false
}

func isCfCommand(s string) bool {
	parts := strings.Split(s, " ")
	for _, cmd := range listCfCommands() {
		if strings.EqualFold(parts[0], cmd) {
			return true
		}
	}
	return false
}

func Executor(s string) {
	s = removeDupSpaces(s)
	if s == "" {
		return
	}
	if isCfCommand(s) {
		executeCfCommand(s)
	} else if isPluginCommand(s) {
		executePluginCommand(s)
	} else {
		executeShellCommand(s)
	}
}

func executeCfCommand(s string) {
	cmdString := "cf"
	args := strings.Split(s, " ")
	cmd := exec.Command(cmdString, args[0:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error during executing %s: %s\n", cmd.Path, err.Error())
	}

	if args[0] == "apps" {
		context.cache.appList = nil
	} else if args[0] == "orgs" {
		context.cache.orgList = nil
	} else if args[0] == "services" {
		context.cache.serviceList = nil
	} else if args[0] == "spaces" {
		context.cache.spaceList = nil
	}
}

func executeShellCommand(s string) {
	cmdString := s
	args := strings.Split(cmdString, " ")

	if strings.ToUpper(args[0]) == "CD" {
		err := os.Chdir(args[1])
		if err != nil {
			fmt.Printf("Error changing CWD to %s: %s\n", args[1], err.Error())
		}
	} else {
		var err error

		if runtime.GOOS == "windows" {
			// not sure which shell so try both cmd and powershell
			for i := 0; i < 2; i++ {
				if i == 0 {
					if isWindowsInternalCommand(s) {
						cmdString = "cmd /c " + s
					}
				} else if i == 1 {
					cmdString = "powershell " + s
				}

				err = execute(cmdString)

				if err == nil {
					break
				}
			}
		} else {
			err = execute(cmdString)
		}

		if err != nil {
			fmt.Printf("Error during executing %s: %s\n", s, err.Error())
		}
	}
}

func execute(cmdString string) error {
	args := strings.Split(cmdString, " ")
	cmd := exec.Command(args[0], args[1:]...)
	dir, _ := os.Getwd()
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
