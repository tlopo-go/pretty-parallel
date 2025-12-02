package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	pr "github.com/tlopo-go/pretty-parallel/parallel-runner"
	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
)

func main() {

	withTerminalCleanup(func() {
		flag.Usage = func() {
			fmt.Println("pretty-parallel [OPTS] < input")
			flag.PrintDefaults()
			notes := regexp.MustCompile(` {16}`).ReplaceAllString(`
                NOTES:
                    Input must be either yaml or json following the schema below:
                    [
                        { name: string, cmd: [string|[]string] },
                        ...
                    ]
            `, "")
			fmt.Println(notes)
		}

		concurrency := flag.Int("c", 10, "concurrency")
		flag.Parse()

		tasks := getTasks()
		runner := pr.New().Concurrency(*concurrency).Tasks(tasks).Run()
		if runner.HasFailed() {
			for _, err := range runner.Errors() {
				log.Error(err)
			}
			os.Exit(1)
		}
	})
}

func withTerminalCleanup(fn func()) {
	// Open controlling terminal
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open /dev/tty: %v\n", err)
		return
	}
	defer tty.Close()
	fd := int(tty.Fd())

	// Backup current state
	oldState, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot get termios: %v\n", err)
		return
	}

	// Ensure state is restored on exit
	defer func() {
		unix.IoctlSetTermios(fd, unix.TIOCSETA, oldState)
	}()
	fn()
}

func getTasks() (tasks []*pr.Task) {
	tasks = []*pr.Task{}
	input := getInput()

	for _, item := range input.([]interface{}) {
		var command []string

		name := item.(map[interface{}]interface{})["name"]
		cmd := item.(map[interface{}]interface{})["cmd"]

		switch v := cmd.(type) {
		case string:
			command = []string{"sh", "-c", v}
		case []string:
			command = v
		default:
			panic("Command must be of type string or []string")
		}

		task := &pr.Task{}
		task.Name = name.(string)
		task.Command = command
		tasks = append(tasks, task)
	}
	return
}

func getInput() (input interface{}) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	input, err = yamlLoad(string(content))
	if err != nil {
		panic(err)
	}

	return
}

func yamlLoad(input string) (data interface{}, err error) {
	err = yaml.Unmarshal([]byte(input), &data)
	return
}
