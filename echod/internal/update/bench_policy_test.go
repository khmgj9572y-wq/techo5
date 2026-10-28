package update

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestBenchMarkerRefusesInstallBeforeValidationOrCallbacks(t *testing.T) {
	old := benchPolicyPath
	benchPolicyPath = filepath.Join(t.TempDir(), "fleet-bench")
	defer func() { benchPolicyPath = old }()
	if err := os.WriteFile(benchPolicyPath, []byte("unsigned native bench only\n"), 0600); err != nil {
		t.Fatal(err)
	}
	err := Install(context.Background(), Manifest{}, func(float32) { t.Fatal("download callback ran") }, func(string) { t.Fatal("install callback ran") })
	if !errors.Is(err, ErrBenchInstallDisabled) {
		t.Fatalf("install was not blocked: %v", err)
	}
	if (Manifest{}).Serves() {
		t.Fatal("bench build advertised an installable update")
	}
}
