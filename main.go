package main

import (
	log "github.com/sirupsen/logrus"
	pr "github.com/tlopo-go/pretty-parallel/parallel-runner"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

func main() {
	tasks := getTasks()
	runner := pr.New().Concurrency(20).Tasks(tasks).Run()
	if runner.HasFailed() {
		for _, err := range runner.Errors() {
			log.Error(err)
		}
	}
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
