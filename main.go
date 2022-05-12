package main

import (
	"fmt"
	"io"
	"os"

	"github.com/onozaty/rcf/replace"
	"github.com/spf13/pflag"
)

const (
	OK    = 0
	Error = 1
)

var (
	Version = "dev"
	Commit  = "none"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {

	var inputPath string
	var outputPath string
	var targetStr string
	var targetRegex string
	var replacement string
	var help bool

	// テストで繰り返しパースすることになるので
	flag := pflag.NewFlagSet("rcf", pflag.ContinueOnError)
	flag.StringVarP(&inputPath, "input", "i", "", "Input file path.")
	flag.StringVarP(&targetRegex, "regex", "r", "", "Target regex.")
	flag.StringVarP(&targetStr, "string", "s", "", "Target string.")
	flag.StringVarP(&replacement, "replacement", "t", "", "Replacement.")
	flag.StringVarP(&outputPath, "output", "o", "", "Output file path.")
	flag.BoolVarP(&help, "help", "h", false, "Help.")
	flag.SortFlags = false
	flag.Usage = func() {
		usage(flag, os.Stderr)
	}

	if err := flag.Parse(args); err != nil {
		usage(flag, os.Stderr)
		fmt.Println("\nError: ", err)
		return Error
	}

	if help {
		usage(flag, os.Stdout)
		return OK
	}

	if inputPath == "" || outputPath == "" || (targetRegex == "" && targetStr == "") {
		usage(flag, os.Stderr)
		return Error
	}

	if err := replaceFiles(inputPath, outputPath, targetRegex, targetStr, replacement); err != nil {
		fmt.Println("\nError: ", err)
		return Error
	}

	return OK
}

func usage(flag *pflag.FlagSet, w io.Writer) {

	fmt.Fprintf(w, "rcf v%s (%s)\n\n", Version, Commit)
	fmt.Fprintf(w, "Usage: rcf -i INPUT [-r REGEX | -s STRING] -t REPLACEMENT -o OUTPUT\n\nFlags\n")
	flag.SetOutput(w)
	flag.PrintDefaults()
}

func replaceFiles(inputPath string, outputPath string, targetRegex string, targetStr string, replacement string) error {

	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var replacer replace.Replacer

	if targetRegex != "" {
		replacer, err = replace.NewRegexpReplacer(targetRegex, replacement)
		if err != nil {
			return err
		}
	} else {
		replacer = replace.NewStringReplacer(targetStr, replacement)
	}

	outputContents := replacer.Replace(string(inputBytes))

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write([]byte(outputContents))
	return err
}
