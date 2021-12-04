package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const OutputEnvName = "OUTPUT"
const BreakLineEnvName = "BREAK_LINE"

func showHelp() {
	fmt.Println("Usage: [OUTPUT=<output log file path|default for stdout>] [BREAK_LINE=<true:default|false>] run-log <commands>")
}

func ParseOutPointer() *os.File {
	out := os.Stdout
	outStr := os.Getenv(OutputEnvName)
	if outStr != "" {
		var err error
		out, err = os.OpenFile(outStr, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			panic(err)
		}
	}
	return out
}

func ParseBreakLine() bool {
	s := os.Getenv(BreakLineEnvName)
	if s == "" {
		return true
	}
	a, _ := strconv.ParseBool(s)
	return a
}

func watchLine(out *os.File, r *bufio.Reader, outType string, w *sync.WaitGroup) {
	outTypeStr := fmt.Sprintf("[%v]", outType)
	line := ""
	for {
		tmp, isPrefix, err := r.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			_, _ = fmt.Fprintln(out, time.Now().Format(time.RFC3339), "[ERR]", err.Error())
		}
		line += string(tmp)
		if isPrefix {
			continue
		}
		_, _ = fmt.Fprintln(out, time.Now().Format(time.RFC3339), outTypeStr, line)
		line = ""
	}
	w.Done()
}

func watchBreakLine(out *os.File, r *bufio.Reader, outType string, w *sync.WaitGroup) {
	outTypeStr := fmt.Sprintf("[%v]", outType)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			_, _ = fmt.Fprintln(out, time.Now().Format(time.RFC3339), "[ERR]", err.Error())
		}

		_, _ = fmt.Fprintln(out, time.Now().Format(time.RFC3339), outTypeStr, string(line))
	}
	w.Done()
}

func main() {
	if len(os.Args) <= 1 || os.Args[1] == "--help" {
		showHelp()
		return
	}

	out := ParseOutPointer()
	defer func() {
		_ = out.Close()
	}()

	breakLine := ParseBreakLine()

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

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Add(1)
	if breakLine {
		go watchBreakLine(out, stdoutPipeReader, "INFO", wg)
		go watchBreakLine(out, stderrPipeReader, "ERR", wg)
	} else {
		go watchLine(out, stdoutPipeReader, "INFO", wg)
		go watchLine(out, stderrPipeReader, "ERR", wg)
	}

	wg.Wait()
}
