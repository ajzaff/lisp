package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/print"
	"github.com/ajzaff/lisp/scan"
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

	var ts scan.TokenScanner
	var s scan.NodeScanner

	var sb strings.Builder
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		switch input := strings.TrimSpace(sc.Text()); {
		case sb.Len() > 0 && input == "":
			ts.Reset(strings.NewReader(sb.String()))
			s.Reset(&ts)
			var nodes []lisp.Node
			for s.Scan() {
				nodes = append(nodes, s.Node())
			}
			if err := s.Err(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				sb.Reset()
				continue
			}
			for _, n := range nodes {
				print.StdPrinter(os.Stdout).Print(n.Val)
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
