package parallelrunner

import (
	"errors"
	"fmt"
	"github.com/tlopo-go/concurrency/executor"
	"github.com/tlopo-go/pretty-parallel/color"
	cr "github.com/tlopo-go/pretty-parallel/command-runner"
	"golang.org/x/term"
	"io/ioutil"
	"os"
	"strings"
)

type Task struct {
	Name    string
	Command []string
}

type ParallelRunner struct {
	tasks       []*Task
	concurrency int
	errors      []error
}

func New() *ParallelRunner {
	return &ParallelRunner{}
}

func (pr *ParallelRunner) Concurrency(value int) *ParallelRunner {
	pr.concurrency = value
	return pr
}

func (pr *ParallelRunner) Tasks(tasks []*Task) *ParallelRunner {
	pr.tasks = tasks
	return pr
}

func (pr *ParallelRunner) Errors() []error {
	return pr.errors
}

func (pr *ParallelRunner) HasFailed() bool {
	return len(pr.errors) > 0
}

func (pr *ParallelRunner) Run() *ParallelRunner {
	tmpDir, err := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	if err != nil {
		panic(err)
	}

	jobs := []executor.Function{}
	pr.errors = []error{}

	longestNameSize := 0
	for _, task := range pr.tasks {
		if len(task.Name) > longestNameSize {
			longestNameSize = len(task.Name)
		}
	}

	for _, task := range pr.tasks {
		padding := longestNameSize - len(task.Name)
		job := func() (err error) {
			err = cr.New().Command(task.Command).Name(task.Name).LogPath(tmpDir).Padding(padding).Run().Error()
			if err != nil {
				err = errors.New(fmt.Sprintf("%s failed, %s", task.Name, err.Error()))
			}
			return
		}
		jobs = append(jobs, job)
	}

	exec := executor.New().Concurrency(pr.concurrency).Jobs(jobs).Run()
	for _, task := range pr.tasks {
		var rightFill string
		fillSize := (getTerminalWidth() - len(task.Name) - 4) / 2
		leftFill := strings.Repeat("=", fillSize)
		if len(task.Name)%2 == 0 {
			rightFill = leftFill
		} else {
			rightFill = leftFill + "="
		}

		fmt.Print(color.Yellow(fmt.Sprintf("%s[ %s ]%s\n", leftFill, task.Name, rightFill)))
		data, _ := os.ReadFile(fmt.Sprintf("%s/.%s.out", tmpDir, task.Name))
		fmt.Print(string(data))
		fmt.Print(color.Yellow(fmt.Sprintf("%s\n", strings.Repeat("=", getTerminalWidth()))))
	}

	if exec.HasFailed() {
		for _, err := range exec.Errors() {
			pr.errors = append(pr.errors, err)
		}
	}
	return pr
}

func getTerminalWidth() int {
	fd := int(os.Stdout.Fd())
	width, _, _ := term.GetSize(fd)

	// if negative it means it's not a TTY, so let's use 160 as default
	if width <= 0 {
		width = 160
	}
	return width
}

func withTempDir(f func(string) error) error {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	return f(tmp)
}
