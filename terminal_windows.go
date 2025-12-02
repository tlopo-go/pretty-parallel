//go:build windows

package main

// withTerminalCleanup on Windows does nothing — Windows does not use /dev/tty or termios.
// It simply runs fn() directly.
func withTerminalCleanup(fn func()) {
	fn()
}
