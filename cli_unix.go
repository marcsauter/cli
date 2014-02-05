package cli

import (
	"os"
	"syscall"
)

const (
	Stty = "/bin/stty"
)

func (c *CLI) EchoOff() error {
	_, err := syscall.ForkExec(Stty, []string{"stty", "-echo"}, &syscall.ProcAttr{Files: []uintptr{os.Stdin.Fd()}})
	if err != nil {
		return err
	}
	return nil
}

func (c CLI) EchoOn() error {
	_, err := syscall.ForkExec(Stty, []string{"stty", "echo"}, &syscall.ProcAttr{Files: []uintptr{os.Stdin.Fd()}})
	if err != nil {
		return err
	}
	return nil
}
