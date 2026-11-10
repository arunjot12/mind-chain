package rustbin

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSelectNewestBinaryPrefersFreshReleaseOverStaleDebug(t *testing.T) {
	targetDir := t.TempDir()
	releaseBin := filepath.Join(targetDir, "release", "mc-reth")
	debugBin := filepath.Join(targetDir, "debug", "mc-reth")
	writeStubBinary(t, debugBin, time.Now().Add(-time.Hour))
	writeStubBinary(t, releaseBin, time.Now())

	got, err := selectNewestBinary(targetDir, "mc-reth")
	if err != nil {
		t.Fatalf("selectNewestBinary: %v", err)
	}
	if got != releaseBin {
		t.Fatalf("expected freshest binary %q, got %q", releaseBin, got)
	}
}

func TestSelectNewestBinaryPrefersFreshDebugOverStaleRelease(t *testing.T) {
	targetDir := t.TempDir()
	releaseBin := filepath.Join(targetDir, "release", "mc-reth")
	debugBin := filepath.Join(targetDir, "debug", "mc-reth")
	writeStubBinary(t, releaseBin, time.Now().Add(-time.Hour))
	writeStubBinary(t, debugBin, time.Now())

	got, err := selectNewestBinary(targetDir, "mc-reth")
	if err != nil {
		t.Fatalf("selectNewestBinary: %v", err)
	}
	if got != debugBin {
		t.Fatalf("expected freshest binary %q, got %q", debugBin, got)
	}
}

func TestSelectNewestBinaryMissing(t *testing.T) {
	if _, err := selectNewestBinary(t.TempDir(), "mc-reth"); err == nil {
		t.Fatal("expected error when no binary is present")
	}
}

func writeStubBinary(t *testing.T, path string, mod time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatal(err)
	}
}
