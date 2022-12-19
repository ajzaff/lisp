//go:build repl
// +build repl

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ajzaff/lisp"
)

const (
	cur  = "> "
	cont = "... "
)

func doRepl() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	go func() {
		s := <-ch
		fmt.Fprintln(os.Stderr, s.String())
		os.Exit(0)
	}()

	sc := bufio.NewScanner(os.Stdin)
	var cons strings.Builder

loop:
	for sc.Scan() {
		input := sc.Text()
		switch {
		case cons.Len() > 0 && input == "":
			no, err := lisp.Parser{}.Parse(cons.String())
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				cons.Reset()
				continue
			}
			cons.Reset()
			lisp.StdPrinter(os.Stdout).Print(no[0].Val)
		case strings.TrimSpace(input) == "quit":
			break loop
		}

		if strings.TrimSpace(input) != "" {
			t := lisp.Tokenizer{}
			if _, err := t.Tokenize(input); err != nil {
				cons.Reset()
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			cons.WriteString(input)
			cons.WriteByte('\n')
		}
	}
}
