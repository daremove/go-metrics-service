package buildversion

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestBuildVersionOutput(t *testing.T) {
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	BuildVersion = "1.0.0"
	BuildDate = "2023-01-01"
	BuildCommit = "abc123"
	printVersion()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = originalStdout

	expectedOutput := "Build version: 1.0.0\nBuild date: 2023-01-01\nBuild commit: abc123\n"
	assert.Equal(t, expectedOutput, buf.String())
}
