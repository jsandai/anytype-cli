package version

import (
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := NewVersionCmd()

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	cmd.SetArgs([]string{"--verbose"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Command with --verbose failed: %v", err)
	}

	verboseFlag := cmd.Flag("verbose")
	if verboseFlag == nil {
		t.Fatal("verbose flag not found")
	}
	if verboseFlag.Shorthand != "v" {
		t.Errorf("verbose flag shorthand = %v, want v", verboseFlag.Shorthand)
	}
}
