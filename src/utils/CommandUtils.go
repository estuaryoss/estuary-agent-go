package utils

import (
	"bytes"
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"github.com/estuaryoss/estuary-agent-go/src/state"
	"log"
	"os/exec"
	"runtime"
)

func RunCommand(command string) *models.CommandDetails {
	cd := models.NewCommandDetails()

	cmd := getOsCommand(command)
	cmd.Env = environment.GetInstance().GetEnvAndVirtualEnvArray()

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
	cd.SetErr(fmt.Sprint(cmd.Stderr))
	cd.SetOut(fmt.Sprint(cmd.Stdout))
	cd.SetCode(cmd.ProcessState.ExitCode())
	cd.SetPid(cmd.ProcessState.Pid())
	cd.SetArgs(cmd.Args)

	log.Printf("Executed command \"%s\", with process id %d\n", command, cmd.ProcessState.Pid())
	return cd
}

func RunCommandToFile(command string, cmdId string) *models.CommandDetails {
	cd := models.NewCommandDetails()

	cmd := getOsCommand(command)
	cmd.Env = environment.GetInstance().GetEnvAndVirtualEnvArray()
	filePathStdOut := GetBase64HashForTheCommand(command, cmdId, ".out")
	filePathStdErr := GetBase64HashForTheCommand(command, cmdId, ".err")
	RecreateFiles([]string{filePathStdOut, filePathStdErr})
	fhStdOut := OpenFile(filePathStdOut)
	fhStdErr := OpenFile(filePathStdErr)
	defer fhStdErr.Close()
	defer fhStdOut.Close()

	cmd.Stdout = fhStdOut
	cmd.Stderr = fhStdErr

	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
	cd.SetErr(string(ReadFile(filePathStdErr)))
	cd.SetOut(string(ReadFile(filePathStdOut)))
	cd.SetCode(cmd.ProcessState.ExitCode())
	cd.SetPid(cmd.ProcessState.Pid())
	cd.SetArgs(cmd.Args)

	log.Printf("Executed command \"%s\", with process id %d\n", command, cmd.ProcessState.Pid())
	return cd
}

func RunCommandNoFile(command string, cmdId string) *models.CommandDetails {
	cd := models.NewCommandDetails()

	cmd := getOsCommand(command)
	cmd.Env = environment.GetInstance().GetEnvAndVirtualEnvArray()

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
	cd.SetErr(fmt.Sprint(cmd.Stderr))
	cd.SetOut(fmt.Sprint(cmd.Stdout))
	cd.SetCode(cmd.ProcessState.ExitCode())
	cd.SetPid(cmd.Process.Pid)
	cd.SetArgs(cmd.Args)

	log.Printf("Executed command \"%s\", with process id %d\n", command, cmd.ProcessState.Pid())
	return cd
}

func StartCommand(cmdId string, command []string, ch chan error) {
	cmd := getCommand(command)
	state.GetInstance().AddCmdToCommandList(cmdId, cmd)

	ch <- cmd.Start()
}

func getCommand(command []string) *exec.Cmd {
	var args []string
	cmd := exec.Command("", args...)

	if runtime.GOOS == "windows" {
		args = []string{"/c"}
		args = append(args, command...)
		cmd = exec.Command("cmd", args...)
	} else {
		args = []string{"-c"}
		args = append(args, command...)
		cmd = exec.Command("sh", args...)
	}

	return cmd
}

/**
! The command must have already ended
*/
func GetCommandDetailsForEndedCommand(endedCommand *exec.Cmd) *models.CommandDetails {
	cd := models.NewCommandDetails()

	cd.SetErr(fmt.Sprint(endedCommand.Stderr))
	cd.SetOut(fmt.Sprint(endedCommand.Stdout))
	cd.SetCode(endedCommand.ProcessState.ExitCode())
	cd.SetPid(endedCommand.ProcessState.Pid())
	cd.SetArgs(endedCommand.Args)

	return cd
}

func getOsCommand(command string) *exec.Cmd {
	var args []string
	cmd := exec.Command("", args...)
	if runtime.GOOS == "windows" {
		args = []string{"/c", command}
		cmd = exec.Command("cmd", args...)
	} else {
		args = []string{"-c", command}
		cmd = exec.Command("sh", args...)
	}
	return cmd
}
