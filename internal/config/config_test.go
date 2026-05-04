package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/devenjarvis/lathe/internal/config"
)

func TestTutorialsDir(t *testing.T) {
	dir, err := config.TutorialsDir()
	if err != nil {
		t.Fatalf("TutorialsDir() error = %v", err)
	}
	if !strings.HasSuffix(dir, ".lathe/tutorials") {
		t.Errorf("TutorialsDir() = %q, want path ending in .lathe/tutorials", dir)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("TutorialsDir() did not create directory at %q", dir)
	}
}
