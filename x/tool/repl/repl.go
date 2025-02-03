package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
	"github.com/ajzaff/lisp/x/print"
)

const (
	cur  = "> "
	cont = "... "
)

func main() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	go func() {
		s := <-ch
		fmt.Fprintln(os.Stderr, s.String())
		os.Exit(0)
	}()

	var ts scan.Scanner
	var s scan.NodeScanner

	var sb strings.Builder
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		switch input := strings.TrimSpace(sc.Text()); {
		case sb.Len() > 0 && input == "":
			ts.Reset(strings.NewReader(sb.String()))
			s.Reset(&ts)
			var vs []lisp.Val
			for s.Scan() {
				_, _, v := s.Node()
				vs = append(vs, v)
			}
			if err := s.Err(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				sb.Reset()
				continue
			}
			for _, v := range vs {
				print.StdPrinter(os.Stdout).Print(v)
			}
			sb.Reset()
		case input == "quit":
			return
		default:
			sb.WriteString(input)
			sb.WriteByte('\n')
		}
	}
}
