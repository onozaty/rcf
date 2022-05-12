package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_Regex(t *testing.T) {

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

func TestRun_String(t *testing.T) {

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
