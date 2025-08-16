package output

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestSuccess(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Success("operation completed %s", "successfully")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Error("Success() should include checkmark")
	}
	if !strings.Contains(output, "operation completed successfully") {
		t.Error("Success() should format message correctly")
	}
}

func TestWarning(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Warning("something might be wrong")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "⚠") {
		t.Error("Warning() should include warning symbol")
	}
	if !strings.Contains(output, "something might be wrong") {
		t.Error("Warning() should include message")
	}
}

func TestError(t *testing.T) {
	err := Error("failed to %s", "connect")

	if err == nil {
		t.Fatal("Error() should return an error")
	}

	if err.Error() != "failed to connect" {
		t.Errorf("Error() = %v, want 'failed to connect'", err.Error())
	}
}

func TestPrint(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print("test %d %s", 123, "message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test 123 message") {
		t.Errorf("Print() output = %v, want 'test 123 message'", output)
	}
}
