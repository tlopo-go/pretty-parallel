package main

import (
	cr "github.com/tlopo-go/pretty-parallel/command-runner"
)

func main() {
	//cmd := []string{"ping", "-c", "8", "8.8.8.8"}
	cmd := []string{"bash", "-c", "seq 1 10 | while read i; do echo $i ; sleep 0.1; done"}
	cr.New().Command(cmd).Name("Job-1").Padding(0).Run()
}
