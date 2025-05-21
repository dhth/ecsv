package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("Skipping integration tests")
	}
}

func TestCLI(t *testing.T) {
	skipIntegration(t)

	tempDir, err := os.MkdirTemp("", "")
	require.NoErrorf(t, err, "error creating temporary directory: %s", err)

	binPath := filepath.Join(tempDir, "ecsv")
	buildArgs := []string{"build", "-o", binPath, "../.."}

	c := exec.Command("go", buildArgs...)
	err = c.Run()
	require.NoErrorf(t, err, "error building binary: %s", err)

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("couldn't clean up temporary directory (%s): %s", binPath, err)
		}
	}()

	// SUCCESSES
	t.Run("Help", func(t *testing.T) {
		// GIVEN
		c := exec.Command(binPath, "-h")

		// WHEN
		b, err := c.CombinedOutput()

		// THEN
		assert.NoError(t, err, "output:\n%s", b)
	})

	t.Run("Parsing correct config works", func(t *testing.T) {
		// GIVEN
		c := exec.Command(
			binPath,
			"--debug",
			"-c",
			"assets/config.yml",
			"-f",
			"html",
		)

		// WHEN
		err := c.Run()

		// THEN
		assert.NoError(t, err)
	})
}
