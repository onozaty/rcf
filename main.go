package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	r "github.com/onozaty/rcf/replace"
	"github.com/spf13/pflag"
)

const (
	OK = 0
	NG = 1
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
	flag.StringVarP(&inputPath, "input", "i", "", "Input file/dir path.")
	flag.StringVarP(&targetRegex, "regex", "r", "", "Target regex.")
	flag.StringVarP(&targetStr, "string", "s", "", "Target string.")
	flag.StringVarP(&replacement, "replacement", "t", "", "Replacement.")
	flag.StringVarP(&outputPath, "output", "o", "", "Output file/dir path.")
	flag.BoolVarP(&help, "help", "h", false, "Help.")
	flag.SortFlags = false
	flag.Usage = func() {
		usage(flag, os.Stderr)
	}

	if err := flag.Parse(args); err != nil {
		usage(flag, os.Stderr)
		fmt.Println("\nError: ", err)
		return NG
	}

	if help {
		usage(flag, os.Stdout)
		return OK
	}

	if inputPath == "" || outputPath == "" || (targetRegex == "" && targetStr == "") {
		usage(flag, os.Stderr)
		return NG
	}

	if err := replace(inputPath, outputPath, targetRegex, targetStr, replacement); err != nil {
		fmt.Println("\nError: ", err)
		return NG
	}

	return OK
}

func usage(flag *pflag.FlagSet, w io.Writer) {

	fmt.Fprintf(w, "rcf v%s (%s)\n\n", Version, Commit)
	fmt.Fprintf(w, "Usage: rcf -i INPUT [-r REGEX | -s STRING] -t REPLACEMENT -o OUTPUT\n\nFlags\n")
	flag.SetOutput(w)
	flag.PrintDefaults()
}

func replace(inputPath string, outputPath string, targetRegex string, targetStr string, replacement string) error {

	replacer, err := newReplacer(targetRegex, targetStr, replacement)
	if err != nil {
		return err
	}

	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	if !inputInfo.IsDir() {
		// ファイル指定
		return replaceFile(inputPath, outputPath, replacer)
	} else {
		// ディレクトリ指定
		return replaceFiles(inputPath, outputPath, replacer)
	}
}

func replaceFiles(inputDirPath string, outputDirPath string, replacer r.Replacer) error {

	entries, err := os.ReadDir(inputDirPath)
	if err != nil {
		return err
	}

	// 出力先のディレクトリが無かったら作っておく
	_, err = os.Stat(outputDirPath)
	if os.IsNotExist(err) {
		os.Mkdir(outputDirPath, os.ModePerm)
	} else if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			err := replaceFile(filepath.Join(inputDirPath, entry.Name()), filepath.Join(outputDirPath, entry.Name()), replacer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func replaceFile(inputFilePath string, outputFilePath string, replacer r.Replacer) error {

	inputBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	outputContents := replacer.Replace(string(inputBytes))

	out, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write([]byte(outputContents))
	return err
}

func newReplacer(targetRegex string, targetStr string, replacement string) (r.Replacer, error) {

	if targetRegex != "" {
		replacer, err := r.NewRegexpReplacer(targetRegex, replacement)
		if err != nil {
			return nil, err
		}
		return replacer, nil
	}

	return r.NewStringReplacer(targetStr, replacement), nil
}
