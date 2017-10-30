package vilka

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var registeredCommands = make(map[string]CommandFunc)

// CommandFunc type for command entrypoint
type CommandFunc func(MsgChan) int

// Register program name and CommandFunc
//
// This should be called in init()
func Register(progName string, cmdFn CommandFunc) {
	registeredCommands[progName] = cmdFn
}

// MaybeExec checks program arguments and if second argument matches registered
// command, it executes it and calls os.Exit() afterwards.
func MaybeExec() {
	if len(os.Args) < 2 {
		return
	}

	cmdFn, ok := registeredCommands[os.Args[1]]
	if !ok {
		return
	}

	msgChan, err := msgChanFromFd(3)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(cmdFn(msgChan))
}

// Launch opens a MsgChan and starts `progName` with other end of MsgChan
func Launch(progName string) (MsgChan, error) {
	if _, ok := registeredCommands[progName]; !ok {
		return nil, fmt.Errorf("Launch: unknown progName: %s", progName)
	}

	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(os.Args[0], progName)
	cmd.ExtraFiles = append(cmd.ExtraFiles, fds[1])
	err = cmd.Start()
	if err != nil {
		// ignore close(2) errors on purpose here
		syscall.Close(fds[0])
		syscall.Close(fds[1])
		return nil, err
	}

	return msgChanFromFd(fds[0])
}
