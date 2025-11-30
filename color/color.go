package color

import "fmt"

func Red(str string) string {
	return fmt.Sprintf("\x1b[31;1m%s\x1b[0m", str)
}

func Green(str string) string {
	return fmt.Sprintf("\x1b[32;1m%s\x1b[0m", str)
}

func Yellow(str string) string {
	return fmt.Sprintf("\x1b[33;1m%s\x1b[0m", str)
}
