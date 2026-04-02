package cleaner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"dupclean/internal/trash"
)

// Protected system paths that should never be deleted
// Note: Path validation is now handled by internal/trash package
var protectedSystemPaths = []string{
	// This is kept for backwards compatibility
	// New code should use trash.MoveToTrash() which has built-in validation
}

// escapeAppleScriptString escapes special characters for AppleScript strings.
// Deprecated: Use internal/trash package instead.
func escapeAppleScriptString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "'", "\\'")
	return s
}

// escapePowerShellString escapes special characters for PowerShell strings.
// Deprecated: Use internal/trash package instead.
func escapePowerShellString(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	return s
}

// validateMediaPath validates that a path is a legitimate media file path.
// Deprecated: Use internal/trash package instead.
func validateMediaPath(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("file does not exist: %w", err)
	}
	return nil
}

// SafePlayMedia plays a media file using OS-native commands with proper escaping.
func SafePlayMedia(path string) (*exec.Cmd, error) {
	if err := validateMediaPath(path); err != nil {
		return nil, err
	}

	absPath, _ := filepath.Abs(path)

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("afplay", absPath), nil
	case "linux":
		return exec.Command("aplay", absPath), nil
	case "windows":
		escapedPath := escapePowerShellString(absPath)
		psScript := fmt.Sprintf(`
$player = New-Object Media.SoundPlayer '%s'
$player.PlaySync()
`, escapedPath)
		return exec.Command("powershell", "-Command", psScript), nil
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// SafeMoveToTrash moves a file to trash using the unified internal/trash package.
// This is a wrapper for backwards compatibility.
func SafeMoveToTrash(path string) error {
	return trash.MoveToTrash(path)
}
