package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	name := os.Args[1]
	args := make([]string, 0)
	for i := 2; i < len(os.Args); i++ {
		args = append(args, os.Args[i])
	}
	cmd := exec.Command(name, args...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stdoutPipeReader := bufio.NewReader(stdoutPipe)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	stderrPipeReader := bufio.NewReader(stderrPipe)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	watchLine := func(r *bufio.Reader, outType string, w *sync.WaitGroup) {
		outTypeStr := fmt.Sprintf("[%v]", outType)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				fmt.Println(time.Now().Format(time.RFC3339), "[ERR]", err.Error())
			}
			fmt.Println(time.Now().Format(time.RFC3339), outTypeStr, string(line))
		}
		w.Done()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Add(1)
	go watchLine(stdoutPipeReader, "INFO", wg)
	go watchLine(stderrPipeReader, "ERR", wg)

	wg.Wait()
}
