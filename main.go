package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	pr "github.com/tlopo-go/pretty-parallel/parallel-runner"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
)

var (
	version string
	commit  string
	date    string
)

func main() {

	withTerminalCleanup(func() {
		flag.Usage = func() {
			fmt.Printf("Version: %s, Commit: %s, Date: %s\n", version, commit, date)
			fmt.Println("\nUSAGE:\npretty-parallel [OPTS] < input")
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
		case []interface{}:
			command = []string{}
			for _, i := range v {
				command = append(command, i.(string))
			}
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
