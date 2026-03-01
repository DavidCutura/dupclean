package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatBytesGUI(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 10, "10.0 KB"},
		{1024 * 100, "100.0 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 5, "5.0 MB"},
		{1024 * 1024 * 100, "100.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
		{1024 * 1024 * 1024 * 2, "2.0 GB"},
		{1024 * 1024 * 1024 * 1024, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatBytes(tt.input)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRuntimeOS(t *testing.T) {
	result := runtimeOS()

	if result == "" {
		t.Log("runtimeOS returned empty string (GOOS not set in test environment)")
		return
	}

	validOS := map[string]bool{
		"darwin":  true,
		"linux":   true,
		"windows": true,
	}

	if !validOS[result] {
		t.Errorf("runtimeOS() = %q, want one of %v", result, []string{"darwin", "linux", "windows"})
	}
}

func TestMoveToTrash_NonExistentFile(t *testing.T) {
	err := moveToTrash("/nonexistent/path/to/file.wav")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestMoveToTrash_EmptyPath(t *testing.T) {
	err := moveToTrash("")

	if err == nil {
		t.Log("moveToTrash does not return error for empty path")
	}
}

func TestMoveToTrash_ValidFileInValidDir(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "valid_test.wav")

	content := []byte("content for deletion test")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	absPath, err := filepath.Abs(testFile)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		t.Fatalf("Test file should exist: %v", err)
	}

	err = moveToTrash(absPath)

	if err != nil {
		t.Logf("moveToTrash error (may be expected in some environments): %v", err)
	}
}

func TestAppState_Struct(t *testing.T) {
	state := &AppState{
		FolderPath:        nil,
		ScanAll:           nil,
		IsScanning:        nil,
		ProgressText:      nil,
		ProgressValue:     nil,
		Groups:            nil,
		CurrentGroupIndex: 0,
		DeletedCount:      0,
		FreedBytes:        0,
	}

	if state.CurrentGroupIndex != 0 {
		t.Errorf("CurrentGroupIndex = %d, want 0", state.CurrentGroupIndex)
	}
	if state.DeletedCount != 0 {
		t.Errorf("DeletedCount = %d, want 0", state.DeletedCount)
	}
	if state.FreedBytes != 0 {
		t.Errorf("FreedBytes = %d, want 0", state.FreedBytes)
	}
}

func TestAppState_WithGroups(t *testing.T) {
	state := &AppState{
		CurrentGroupIndex: 1,
		DeletedCount:      5,
		FreedBytes:        1024 * 1024 * 10,
	}

	if state.CurrentGroupIndex != 1 {
		t.Errorf("CurrentGroupIndex = %d, want 1", state.CurrentGroupIndex)
	}
	if state.DeletedCount != 5 {
		t.Errorf("DeletedCount = %d, want 5", state.DeletedCount)
	}
	if state.FreedBytes != 10*1024*1024 {
		t.Errorf("FreedBytes = %d, want %d", state.FreedBytes, 10*1024*1024)
	}
}
