//go:build linux
// +build linux

package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

func withTerminalCleanup(fn func()) {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fn()
		return
	}
	defer tty.Close()
	fd := int(tty.Fd())

	oldState, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get termios: %v\n", err)
		fn()
		return
	}

	defer unix.IoctlSetTermios(fd, unix.TCSETS, oldState)

	fn()
}
