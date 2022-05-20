package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

func TestRun_File_Regex(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "abc\nabc\naa")
	output := filepath.Join(d, "output.txt")

	args := []string{
		"-i", input,
		"-r", "a$",
		"-t", "x",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, output)
	assert.Equal(t, "abc\nabc\nax", replaced)
}

func TestRun_File_String(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "aa.ab.ac.ad.a.b.c.d")
	output := filepath.Join(d, "output.txt")

	args := []string{
		"-i", input,
		"-s", "a.",
		"-t", "xx",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, output)
	assert.Equal(t, "axxab.ac.ad.xxb.c.d", replaced)
}

func TestRun_Dir_Regex(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createDir(t, d, "input")

	createFileWriteString(t, input, "input1.txt", "abc\nabc\naa")
	createFileWriteString(t, input, "input2.txt", "a")
	createFileWriteString(t, input, "input3.txt", "ax")

	output := createDir(t, d, "output")

	args := []string{
		"-i", input,
		"-r", "a$",
		"-t", "x",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	{
		replaced := readString(t, filepath.Join(output, "input1.txt"))
		assert.Equal(t, "abc\nabc\nax", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "input2.txt"))
		assert.Equal(t, "x", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "input3.txt"))
		assert.Equal(t, "ax", replaced)
	}
}

func TestRun_File_Regex_Japanese(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "あいうえおかきくけこ")
	output := filepath.Join(d, "output.txt")

	args := []string{
		"-i", input,
		"-r", "あ.{4}",
		"-t", "",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, output)
	assert.Equal(t, "かきくけこ", replaced)
}

func TestRun_Dir_String(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createDir(t, d, "input")

	createFileWriteString(t, input, "input1.txt", "abc\na.c\naa")
	createFileWriteString(t, input, "input2.txt", "")
	createFileWriteString(t, input, "input3.txt", "a.c")

	output := createDir(t, d, "output")

	args := []string{
		"-i", input,
		"-s", "a.c",
		"-t", "",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	{
		replaced := readString(t, filepath.Join(output, "input1.txt"))
		assert.Equal(t, "abc\n\naa", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "input2.txt"))
		assert.Equal(t, "", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "input3.txt"))
		assert.Equal(t, "", replaced)
	}
}

func TestRun_Dir_CreateOutputDir(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createDir(t, d, "input")

	createFileWriteString(t, input, "input1.txt", "abc\nabc\nabc")

	output := filepath.Join(d, "output") // 出力ディレクトリは存在しない状態

	args := []string{
		"-i", input,
		"-r", "(?m)c$",
		"-t", "x",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, filepath.Join(output, "input1.txt"))
	assert.Equal(t, "abx\nabx\nabx", replaced)
}

func TestRun_Dir_Recursive(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createDir(t, d, "input")
	createFileWriteString(t, input, "1.txt", "abc")
	createFileWriteString(t, input, "2.txt", "")

	inputSub := createDir(t, input, "sub")
	createFileWriteString(t, inputSub, "3.txt", "cat")
	createFileWriteString(t, inputSub, "4.txt", "aaa")

	inputSubSub := createDir(t, inputSub, "sub")
	createFileWriteString(t, inputSubSub, "5.txt", "a")

	output := createDir(t, d, "output")

	args := []string{
		"-i", input,
		"-s", "a",
		"-t", "",
		"-R",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	{
		replaced := readString(t, filepath.Join(output, "1.txt"))
		assert.Equal(t, "bc", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "2.txt"))
		assert.Equal(t, "", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "sub", "3.txt"))
		assert.Equal(t, "ct", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "sub", "4.txt"))
		assert.Equal(t, "", replaced)
	}
	{
		replaced := readString(t, filepath.Join(output, "sub", "sub", "5.txt"))
		assert.Equal(t, "", replaced)
	}
}

func TestRun_Escape_String(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "1\n2\n")
	output := filepath.Join(d, "output.txt")

	args := []string{
		"-i", input,
		"-s", `\n`,
		"-t", `\t`,
		"--escape",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, output)
	assert.Equal(t, "1\t2\t", replaced)
}

func TestRun_Escape_Regex(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "a　　　")
	output := filepath.Join(d, "output.txt")

	args := []string{
		"-i", input,
		"-r", `\u3000+`,
		"-t", `\u0020`,
		"-e",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, output)
	assert.Equal(t, "a ", replaced)
}

func TestRun_Overwrite(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "aaa")

	args := []string{
		"-i", input,
		"-s", "a",
		"-t", "x",
		"--overwrite",
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, input)
	assert.Equal(t, "xxx", replaced)
}

func TestRun_Charset_UTF8(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "あいうえお")
	output := filepath.Join(d, "output.txt")

	args := []string{
		"-i", input,
		"-r", "あ.う",
		"-t", "",
		"-c", "utf-8",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := readString(t, output)
	assert.Equal(t, "えお", replaced)
}

func TestRun_Charset_SJIS(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteBytes(t, d, "input.txt", stringToByte(t, "あいうえお", japanese.ShiftJIS))
	output := filepath.Join(d, "output.txt")

	args := []string{
		"-i", input,
		"-r", "あ.う",
		"-t", "",
		"-c", "sjis",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)
	replaced := byteToString(t, readBytes(t, output), japanese.ShiftJIS)
	assert.Equal(t, "えお", replaced)
}

func TestRun_Charset_Invalid(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "")
	output := filepath.Join(d, "output.txt")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", input,
		"-r", "x",
		"-t", "",
		"-c", "xxxx", // 存在しないCharset
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "\nError: htmlindex: invalid encoding name\n", buf.String())
}

func TestRun_InvalidRegex(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "")
	output := filepath.Join(d, "output.txt")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", input,
		"-r", "[a", // 不正な正規表現
		"-t", "",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "\nError: error parsing regexp: missing closing ]: `[a`\n", buf.String())
}

func TestRun_InvalidEscape_Regex(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "")
	output := filepath.Join(d, "output.txt")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", input,
		"-r", `\x`, // 不正なエスケープ
		"-t", "",
		"-e",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "\nError: --regex is invalid string: \\x\n", buf.String())
}

func TestRun_InvalidEscape_String(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "")
	output := filepath.Join(d, "output.txt")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", input,
		"-s", `\x`, // 不正なエスケープ
		"-t", "",
		"-e",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "\nError: --string is invalid string: \\x\n", buf.String())
}

func TestRun_InvalidEscape_Replacement(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createFileWriteString(t, d, "input.txt", "")
	output := filepath.Join(d, "output.txt")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", input,
		"-s", "a",
		"-t", `\`, // 不正なエスケープ
		"-e",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Equal(t, "\nError: --replacement is invalid string: \\\n", buf.String())
}

func TestRun_InputNotFound(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := filepath.Join(d, "input") // 存在しない
	output := filepath.Join(d, "output")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", input,
		"-s", "a",
		"-t", "",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "input: The system cannot find the file specified")
}

func TestRun_OutputNotFound(t *testing.T) {

	// ARRANGE
	d := createTempDir(t)
	defer os.RemoveAll(d)

	input := createDir(t, d, "input")
	output := filepath.Join(d, "a", "b") // 親ディレクトリ自体が無い

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", input,
		"-s", "a",
		"-t", "",
		"-o", output,
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "b: The system cannot find the path specified.")
}

func TestRun_Help(t *testing.T) {

	// ARRANGE
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	args := []string{
		"-h",
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, OK, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "Usage: rcf")
}

func TestRun_InvalidArgs_Empty(t *testing.T) {

	// ARRANGE
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "Usage: rcf")
}

func TestRun_InvalidArgs_Unknown(t *testing.T) {

	// ARRANGE
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-x",
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "unknown shorthand flag: 'x' in -x")
}

func TestRun_InvalidArgs_InputEmpty(t *testing.T) {

	// ARRANGE
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-s", "a",
		"-t", "",
		"-o", "out",
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "Usage: rcf")
}

func TestRun_InvalidArgs_StringAndRegexEmpty(t *testing.T) {

	// ARRANGE
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", "in",
		"-t", "",
		"-o", "out",
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "Usage: rcf")
}

func TestRun_InvalidArgs_OutputAndOverwriteEmpty(t *testing.T) {

	// ARRANGE
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = stderr }()

	args := []string{
		"-i", "in",
		"-s", "x",
		"-t", "",
	}

	// ACT
	c := run(args)

	// ASSERT
	require.Equal(t, NG, c)

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	assert.Contains(t, buf.String(), "Usage: rcf")
}

//////////////////////////////////////////////////////////

func createFileWriteBytes(t *testing.T, dir string, name string, content []byte) string {

	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.Write(content)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	return file.Name()
}

func createFileWriteString(t *testing.T, dir string, name string, content string) string {

	return createFileWriteBytes(t, dir, name, []byte(content))
}

func createTempDir(t *testing.T) string {

	tempDir, err := os.MkdirTemp("", "rcf")
	if err != nil {
		t.Fatal(err)
	}

	return tempDir
}

func createDir(t *testing.T, parent string, name string) string {

	dir := filepath.Join(parent, name)
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func readBytes(t *testing.T, name string) []byte {

	bo, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	return bo
}

func readString(t *testing.T, name string) string {

	bo := readBytes(t, name)
	return string(bo)
}

func byteToString(t *testing.T, bytes []byte, enc encoding.Encoding) string {
	decoded, err := enc.NewDecoder().Bytes(bytes)
	if err != nil {
		t.Fatal(err)
	}

	return string(decoded)
}

func stringToByte(t *testing.T, str string, enc encoding.Encoding) []byte {
	encoded, err := enc.NewEncoder().Bytes([]byte(str))
	if err != nil {
		t.Fatal(err)
	}

	return encoded
}
