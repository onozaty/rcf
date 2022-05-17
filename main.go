package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	r "github.com/onozaty/rcf/replace"
	"github.com/spf13/pflag"
	e "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
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
	var escapeSequence bool
	var charset string
	var overwrite bool
	var help bool

	// テストで繰り返しパースすることになるので
	flag := pflag.NewFlagSet("rcf", pflag.ContinueOnError)
	flag.StringVarP(&inputPath, "input", "i", "", "Input file/dir path.")
	flag.StringVarP(&targetRegex, "regex", "r", "", "Target regex.")
	flag.StringVarP(&targetStr, "string", "s", "", "Target string.")
	flag.StringVarP(&replacement, "replacement", "t", "", "Replacement.")
	flag.BoolVarP(&escapeSequence, "escape", "e", false, "Enable escape sequence.")
	flag.StringVarP(&charset, "charset", "c", "UTF-8", "Charset. (default UTF-8)")
	flag.StringVarP(&outputPath, "output", "o", "", "Output file/dir path.")
	flag.BoolVarP(&overwrite, "overwrite", "O", false, "Overwrite the input file.")
	flag.BoolVarP(&help, "help", "h", false, "Help.")
	flag.SortFlags = false
	flag.Usage = func() {
		usage(flag, os.Stderr)
	}

	if err := flag.Parse(args); err != nil {
		usage(flag, os.Stderr)
		fmt.Fprintln(os.Stderr, "\nError:", err)
		return NG
	}

	if help {
		usage(flag, os.Stdout)
		return OK
	}

	if inputPath == "" || (outputPath == "" && !overwrite) || (targetRegex == "" && targetStr == "") {
		usage(flag, os.Stderr)
		return NG
	}

	if escapeSequence {
		// Unquoteした文字列を再設定
		if unquoted, err := unquote(targetRegex); err != nil {
			fmt.Fprintln(os.Stderr, "\nError: --regex is invalid string:", targetRegex)
			return NG
		} else {
			targetRegex = unquoted
		}

		if unquoted, err := unquote(targetStr); err != nil {
			fmt.Fprintln(os.Stderr, "\nError: --string is invalid string:", targetStr)
			return NG
		} else {
			targetStr = unquoted
		}

		if unquoted, err := unquote(replacement); err != nil {
			fmt.Fprintln(os.Stderr, "\nError: --replacement is invalid string:", replacement)
			return NG
		} else {
			replacement = unquoted
		}
	}

	if outputPath == "" && overwrite {
		// 上書き指定されていた場合、入力と同じものを指定
		outputPath = inputPath
	}

	condition := condition{
		targetRegex: targetRegex,
		targetStr:   targetStr,
		replacement: replacement,
	}

	if err := replace(inputPath, outputPath, condition, charset); err != nil {
		fmt.Fprintln(os.Stderr, "\nError:", err)
		return NG
	}

	return OK
}

func usage(flag *pflag.FlagSet, w io.Writer) {

	fmt.Fprintf(w, "rcf v%s (%s)\n\n", Version, Commit)
	fmt.Fprintf(w, "Usage: rcf -i INPUT [-r REGEX | -s STRING] -t REPLACEMENT [--escape] [-c CHARSET] [-o OUTPUT | --overwrite]\n\nFlags\n")
	flag.SetOutput(w)
	flag.PrintDefaults()
}

type condition struct {
	targetRegex string
	targetStr   string
	replacement string
}

func replace(inputPath string, outputPath string, condition condition, charset string) error {

	encoding, err := htmlindex.Get(charset)
	if err != nil {
		return err
	}

	replacer, err := newReplacer(condition)
	if err != nil {
		return err
	}

	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	if !inputInfo.IsDir() {
		// ファイル指定
		return replaceFile(inputPath, outputPath, replacer, encoding)
	} else {
		// ディレクトリ指定
		return replaceFiles(inputPath, outputPath, replacer, encoding)
	}
}

func replaceFiles(inputDirPath string, outputDirPath string, replacer r.Replacer, encoding e.Encoding) error {

	entries, err := os.ReadDir(inputDirPath)
	if err != nil {
		return err
	}

	// 出力先のディレクトリが無かったら作っておく
	_, err = os.Stat(outputDirPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(outputDirPath, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			err := replaceFile(filepath.Join(inputDirPath, entry.Name()), filepath.Join(outputDirPath, entry.Name()), replacer, encoding)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func replaceFile(inputFilePath string, outputFilePath string, replacer r.Replacer, encoding e.Encoding) error {

	inputBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	decodedBytes, err := encoding.NewDecoder().Bytes(inputBytes)
	if err != nil {
		return err
	}

	outputContents := replacer.Replace(string(decodedBytes))

	out, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	encodedBytes, err := encoding.NewEncoder().Bytes([]byte(outputContents))
	if err != nil {
		return err
	}

	_, err = out.Write(encodedBytes)
	return err
}

func newReplacer(condition condition) (r.Replacer, error) {

	if condition.targetRegex != "" {
		replacer, err := r.NewRegexpReplacer(condition.targetRegex, condition.replacement)
		if err != nil {
			return nil, err
		}
		return replacer, nil
	}

	return r.NewStringReplacer(condition.targetStr, condition.replacement), nil
}

func unquote(str string) (string, error) {
	return strconv.Unquote(`"` + str + `"`)
}
