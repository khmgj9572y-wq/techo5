package echod

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HuskerMinion/techo5/echod/internal/config"
	"github.com/HuskerMinion/techo5/echod/internal/feature/detect/assets"
)

func TestWakeWAVUsesTheSavedThresholdInsteadOfZero(t *testing.T) {
	dir := t.TempDir()
	config.Use(filepath.Join(dir, "state.json"))
	models := filepath.Join(dir, "models")
	if _, err := assets.Stop(models); err != nil {
		t.Fatal(err)
	}
	// Four seconds of canonical 16 kHz mono S16 silence.
	wav := make([]byte, 44+4*16000*2)
	copy(wav, "RIFF")
	binary.LittleEndian.PutUint32(wav[4:], uint32(len(wav)-8))
	copy(wav[8:], "WAVEfmt ")
	binary.LittleEndian.PutUint32(wav[16:], 16)
	binary.LittleEndian.PutUint16(wav[20:], 1)
	binary.LittleEndian.PutUint16(wav[22:], 1)
	binary.LittleEndian.PutUint32(wav[24:], 16000)
	binary.LittleEndian.PutUint32(wav[28:], 32000)
	binary.LittleEndian.PutUint16(wav[32:], 2)
	binary.LittleEndian.PutUint16(wav[34:], 16)
	copy(wav[36:], "data")
	binary.LittleEndian.PutUint32(wav[40:], uint32(len(wav)-44))
	path := filepath.Join(dir, "silence.wav")
	if err := os.WriteFile(path, wav, 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd := newWakeCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--models", models, "--wav", path})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "DETECTED") || !strings.Contains(out.String(), "200 frames, 0 detections") {
		t.Fatalf("silence was treated as a wake word: %s", &out)
	}
}
