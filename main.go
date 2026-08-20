package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"expr-eval/internal/eval"
	"expr-eval/internal/parser"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: expr-eval <expression> [-name value]...")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}
	if err := run(args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	expr, vars, err := parseArgs(args)
	if err != nil {
		return err
	}
	node, err := parser.Parse(expr)
	if err != nil {
		return err
	}
	result, err := eval.Eval(node, vars)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func parseArgs(args []string) (string, map[string]float64, error) {
	vars := make(map[string]float64)
	var expr string
	exprSet := false
	var rest []string
	for _, a := range args {
		if !exprSet && !strings.HasPrefix(a, "-") {
			expr = a
			exprSet = true
			continue
		}
		rest = append(rest, a)
	}
	if !exprSet {
		return "", nil, fmt.Errorf("missing expression")
	}
	for i := 0; i < len(rest); i++ {
		a := rest[i]
		if !strings.HasPrefix(a, "-") {
			return "", nil, fmt.Errorf("unexpected argument: %s", a)
		}
		name := a[1:]
		if eq := strings.IndexByte(name, '='); eq >= 0 {
			v, e := strconv.ParseFloat(name[eq+1:], 64)
			if e != nil {
				return "", nil, fmt.Errorf("invalid value for %s: %v", name[:eq], e)
			}
			vars[name[:eq]] = v
			continue
		}
		if i+1 >= len(rest) {
			return "", nil, fmt.Errorf("missing value for %s", name)
		}
		i++
		v, e := strconv.ParseFloat(rest[i], 64)
		if e != nil {
			return "", nil, fmt.Errorf("invalid value for %s: %v", name, e)
		}
		vars[name] = v
	}
	return expr, vars, nil
}
