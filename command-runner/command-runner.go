package commandrunner

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tlopo-go/pretty-parallel/color"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CommandRunner struct {
	name              string
	command           []string
	padding           int
	err               error
	outputFile        string
	outputFileHandler *os.File
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

func (cr *CommandRunner) Error() error {
	return cr.err
}

func (cr *CommandRunner) Run() *CommandRunner {
	var err error
	outPrefix := color.Green(fmt.Sprintf("%s%s |", strings.Repeat(" ", cr.padding), cr.name))
	errPrefix := color.Red(fmt.Sprintf("%s%s |", strings.Repeat(" ", cr.padding), cr.name))

	cr.outputFile = fmt.Sprintf(".%s.out", cr.name)
	cr.outputFileHandler, err = os.OpenFile(cr.outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer cr.outputFileHandler.Close()

	writer := io.MultiWriter(os.Stdout, cr.outputFileHandler)

	cmd := exec.Command(cr.command[0], cr.command[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		cr.err = err
		return cr
	}

	// Stream stdout
	go stream(outPrefix, writer, stdout)

	// Stream stderr
	go stream(errPrefix, writer, stderr)

	// Wait for the command to exit
	if err := cmd.Wait(); err != nil {
		cr.err = errors.New(fmt.Sprintf("Command exited with error: %+v", err))
		logErrorToOutputFile(errPrefix, writer, cr.err)
	}
	return cr
}

func logErrorToOutputFile(prefix string, w io.Writer, err error) {
	for _, line := range strings.Split(fmt.Sprintf("%+v", err), "\n") {
		fmt.Fprintf(w, "%s %s\n", prefix, line)
	}
}

func stream(prefix string, w io.Writer, r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(prefix, "read error:", err)
			}
			break
		}
		fmt.Fprintf(w, "%s %s", prefix, line)
	}
}
