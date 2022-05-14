package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
	r := run(args)
	if r != OK {
		t.Fatal("run failed")
	}

	// ASSERT
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
	r := run(args)
	if r != OK {
		t.Fatal("run failed")
	}

	// ASSERT
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
	r := run(args)
	if r != OK {
		t.Fatal("run failed")
	}

	// ASSERT
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
	r := run(args)
	if r != OK {
		t.Fatal("run failed")
	}

	// ASSERT
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
