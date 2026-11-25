package homedir

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDirUnixGetentNotFound(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "plan9" {
		t.Skip("the getent lookup is only used on non-darwin, non-plan9 Unix systems")
	}

	// Look up sh before PATH is replaced below.
	shPath, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not available")
	}
	shPath, err = filepath.EvalSymlinks(shPath)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Build a PATH that contains a working sh but no getent. Running getent
	// then fails with an *exec.Error wrapping exec.ErrNotFound, and dirUnix
	// must ignore that error and fall through to the shell fallback.
	tmp := t.TempDir()
	if err := os.Symlink(shPath, filepath.Join(tmp, "sh")); err != nil {
		t.Fatalf("err: %s", err)
	}

	defer patchEnv("PATH", tmp)()
	defer patchEnv("HOME", "")()

	result, err := dirUnix()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if result == "" {
		t.Fatal("expected a home directory from the shell fallback, got empty string")
	}
}
