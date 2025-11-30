package commandrunner

import (
	"bufio"
	"fmt"
	"github.com/tlopo-go/pretty-parallel/color"
	"io"
	"os/exec"
	"strings"
)

type CommandRunner struct {
	name    string
	command []string
	padding int
}

func New() (cr *CommandRunner) {
	cr = &CommandRunner{}
	cr.name = "no-name"
	cr.padding = 0
	return
}

func (cr *CommandRunner) Name(name string) *CommandRunner {
	cr.name = name
	return cr
}

func (cr *CommandRunner) Command(command []string) *CommandRunner {
	cr.command = command
	return cr
}

func (cr *CommandRunner) Padding(padding int) *CommandRunner {
	cr.padding = padding
	return cr
}

func (cr *CommandRunner) Run() *CommandRunner {
	outPrefix := color.Green(fmt.Sprintf("%s%s |", strings.Repeat(" ", cr.padding), cr.name))
	errPrefix := color.Red(fmt.Sprintf("%s%s |", strings.Repeat(" ", cr.padding), cr.name))

	cmd := exec.Command(cr.command[0], cr.command[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// Stream stdout
	go stream(outPrefix, stdout)

	// Stream stderr
	go stream(errPrefix, stderr)

	// Wait for the command to exit
	if err := cmd.Wait(); err != nil {
		fmt.Println("Command exited with error:", err)
	}
	return cr
}

func stream(prefix string, r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(prefix, "read error:", err)
			}
			break
		}
		fmt.Printf("%s %s", prefix, line)
	}
}
