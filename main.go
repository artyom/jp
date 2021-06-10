// TODO describe program
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jmespath/go-jmespath"
)

func main() {
	log.SetFlags(0)
	args := runArgs{}
	flag.BoolVar(&args.pretty, "p", args.pretty, "format output nicely")
	flag.Parse()
	args.query = flag.Arg(0)
	args.input = flag.Arg(1)
	if err := run(args); err != nil {
		if err == errInvalidUsage {
			flag.Usage()
			os.Exit(2)
		}
		log.Fatal(err)
	}
}

type runArgs struct {
	pretty bool
	query  string
	input  string
}

func run(args runArgs) error {
	if args.query == "" || args.input == "" {
		return errInvalidUsage
	}
	jp, err := jmespath.Compile(args.query)
	if err != nil {
		return fmt.Errorf("compiling query %q: %w", args.query, err)
	}
	b, err := os.ReadFile(args.input)
	if err != nil {
		return err
	}
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	out, err := jp.Search(data)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if args.pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(out)
}

var errInvalidUsage = errors.New("both query and input must be set")

const usage = `Usage: jp [flags] <query> <file>`

func init() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "For query syntax see https://jmespath.org")
	}
}
