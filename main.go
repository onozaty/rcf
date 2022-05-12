package main

import (
	"fmt"
	"io"
	"os"

	"github.com/onozaty/rcf/replace"
	flag "github.com/spf13/pflag"
)

var (
	Version = "dev"
	Commit  = "none"
)

func main() {

	var inputPath string
	var outputPath string
	var targetStr string
	var targetRegex string
	var replacement string
	var help bool

	flag.StringVarP(&inputPath, "input", "i", "", "Input file path.")
	flag.StringVarP(&targetRegex, "regex", "r", "", "Target regex.")
	flag.StringVarP(&targetStr, "string", "s", "", "Target string.")
	flag.StringVarP(&replacement, "replacement", "t", "", "Replacement.")
	flag.StringVarP(&outputPath, "output", "o", "", "Output file path.")
	flag.BoolVarP(&help, "help", "h", false, "Help.")
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		usage(os.Stderr)
	}

	flag.Parse()

	if help {
		usage(os.Stdout)
		os.Exit(0)
	}

	if inputPath == "" || outputPath == "" || (targetRegex == "" && targetStr == "") {
		usage(os.Stderr)
		os.Exit(1)
	}

	err := run(inputPath, outputPath, targetRegex, targetStr, replacement)
	if err != nil {
		fmt.Println("\nError: ", err)
		os.Exit(1)
	}
}

func usage(w io.Writer) {

	fmt.Fprintf(w, "rcf v%s (%s)\n\n", Version, Commit)
	fmt.Fprintf(w, "Usage: rcf -i INPUT [-r REGEX | -s STRING] -t REPLACEMENT -o OUTPUT\n\nFlags\n")
	flag.CommandLine.SetOutput(w)
	flag.PrintDefaults()
}

func run(inputPath string, outputPath string, targetRegex string, targetStr string, replacement string) error {

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
