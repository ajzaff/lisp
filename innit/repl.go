package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ajzaff/innit"
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
	var expr strings.Builder

loop:
	for sc.Scan() {
		input := sc.Text()
		switch {
		case expr.Len() > 0 && input == "":
			no, err := innit.Parse(expr.String())
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				expr.Reset()
				continue
			}
			expr.Reset()
			innit.CompactPrinter(os.Stdout).Print(no)
		case strings.TrimSpace(input) == "quit":
			break loop
		}

		if strings.TrimSpace(input) != "" {
			if _, err := innit.Tokenize(input); err != nil {
				expr.Reset()
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			expr.WriteString(input)
			expr.WriteByte('\n')
		}
	}
}
