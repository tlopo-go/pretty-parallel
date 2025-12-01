package commandrunner

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/tlopo-go/concurrency/thread"
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
	logPath           string
	outputFileHandler *os.File
}

func New() (cr *CommandRunner) {
	cr = &CommandRunner{}
	cr.name = "no-name"
	cr.padding = 0
	cr.logPath = "/tmp"
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

func (cr *CommandRunner) LogPath(value string) *CommandRunner {
	cr.logPath = value
	return cr
}

func (cr *CommandRunner) Error() error {
	return cr.err
}

func (cr *CommandRunner) Run() *CommandRunner {
	var err error
	outPrefix := color.Green(fmt.Sprintf("%s%s |", strings.Repeat(" ", cr.padding), cr.name))
	errPrefix := color.Red(fmt.Sprintf("%s%s |", strings.Repeat(" ", cr.padding), cr.name))

	cr.outputFile = fmt.Sprintf("%s/.%s.out", cr.logPath, cr.name)
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
	outThread := thread.New(func() (err error) {
		stream(outPrefix, writer, stdout)
		return
	}).Start()

	// Stream stderr
	errThread := thread.New(func() (err error) {
		stream(errPrefix, writer, stderr)
		return
	}).Start()

	outThread.Join()
	errThread.Join()

	// Wait for the command to exit
	if err := cmd.Wait(); err != nil {
		cr.err = errors.New(fmt.Sprintf("Command exited with error: %+v", err))
		logErrorToOutputFile(errPrefix, writer, cr.err)
	}
	return cr
}

func logErrorToOutputFile(prefix string, w io.Writer, err error) {
	for _, line := range strings.Split(err.Error(), "\n") {
		fmt.Fprintf(w, "%s %s\n", prefix, line)
	}
}

func stream(prefix string, w io.Writer, r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			if err != io.EOF {
				fmt.Println(prefix, "read error:", err)
			}
			break
		}

		if len(line) > 0 {
			fmt.Fprintf(w, "%s %s\n", prefix, strings.TrimRight(line, "\n"))
		}

		if err == io.EOF {
			break
		}

	}
}
