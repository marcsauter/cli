package cli

import (
	"syscall"
	"unsafe"
)

const (
	ENABLE_ECHO_INPUT = uint32(0x0004)
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode = modkernel32.NewProc("GetConsoleMode")
	setConsoleMode = modkernel32.NewProc("SetConsoleMode")
)

func (c *CLI) EchoOff() error {
	var h syscall.Handle
	h, err := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall(getConsoleMode.Addr(), 2, uintptr(h), uintptr(unsafe.Pointer(&c.consoleMode)), 0)
	if errno != 0 {
		return errno
	}
	_, _, errno = syscall.Syscall(setConsoleMode.Addr(), 2, uintptr(h), uintptr(c.consoleMode & ^(ENABLE_ECHO_INPUT)), 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func (c CLI) EchoOn() error {
	h, err := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall(setConsoleMode.Addr(), 2, uintptr(h), uintptr(c.consoleMode), 0)
	if errno != 0 {
		return errno
	}
	return nil
}
