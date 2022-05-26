package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return 1
	}
	name, args := cmd[0], cmd[1:]
	comm := exec.Command(name, args...)
	envsToSet, envsToDelete := env.GetEnvs()
	for k, v := range envsToSet {
		if err := os.Setenv(k, v); err != nil {
			return 1
		}
	}
	for _, v := range envsToDelete {
		if _, ok := os.LookupEnv(v); ok {
			if err := os.Unsetenv(v); err != nil {
				return 1
			}
		}
	}
	comm.Stdout, comm.Stdin, comm.Stderr = os.Stdout, os.Stdin, os.Stderr
	if err := comm.Run(); err != nil {
		return 1
	}
	return 0
}
