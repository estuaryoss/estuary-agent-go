package utils

import (
	"estuary-agent-go/src/environment"
	"estuary-agent-go/src/models"
	"os/exec"
	"runtime"
)

func RunCommand(command string) *models.CommandDetails {
	cd := models.NewCommandDetails()
	env := environment.GetInstance()

	cmd := exec.Command("", command)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", []string{"/c", command}...)
	} else {
		cmd = exec.Command("sh", []string{"-c", command}...)
	}
	cmd.Env = env.GetEnvAndVirtualEnvArray()

	err := cmd.Run()
	if err != nil {
		cd.SetErr(err.Error())
	}

	cd.SetOut(cmd.ProcessState.String())
	cd.SetCode(cmd.ProcessState.ExitCode())
	cd.SetPid(100)
	cd.SetArgs(command)

	//log.Printf("Executed command %s, with process id %d\n", command, cmd.ProcessState.Pid())
	return cd
}
