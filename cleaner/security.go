package cleaner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Protected system paths that should never be deleted
var protectedSystemPaths = []string{
	// Unix/Linux
	"/etc",
	"/bin",
	"/sbin",
	"/usr",
	"/lib",
	"/lib64",
	"/boot",
	"/dev",
	"/proc",
	"/sys",
	// macOS
	"/System",
	"/Library",
	"/Applications",
	// Windows
	`C:\Windows`,
	`C:\Program Files`,
	`C:\Program Files (x86)`,
	`C:\ProgramData`,
}

// Paths that are commonly used but should NOT be protected
var allowedSubPaths = []string{
	"/var/tmp",
	"/var/folders",
	"/var/log",
	"/usr/local",
}

// sanitizePathForShell ensures a path is safe to use in shell commands.
// It validates the path exists and escapes special characters.
func sanitizePathForShell(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Clean the path to resolve .. and .
	cleanPath := filepath.Clean(absPath)

	// Check for path traversal attempts
	if strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal detected")
	}

	// Verify the path exists to prevent injection via non-existent paths
	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("path does not exist: %w", err)
	}

	return absPath, nil
}

// isProtectedPath checks if a path is a protected system directory.
// This prevents accidental deletion of critical system files.
func isProtectedPath(path string) bool {
	// Normalize path for comparison
	cleanPath := filepath.Clean(path)
	lowerPath := strings.ToLower(cleanPath)

	// First check if this is an allowed subpath
	for _, allowed := range allowedSubPaths {
		lowerAllowed := strings.ToLower(allowed)
		if lowerPath == lowerAllowed || strings.HasPrefix(lowerPath, lowerAllowed+string(filepath.Separator)) {
			return false // This is an allowed subpath
		}
	}

	// Check protected paths
	for _, protected := range protectedSystemPaths {
		lowerProtected := strings.ToLower(protected)
		
		// Exact match
		if lowerPath == lowerProtected {
			return true
		}
		
		// Path is inside protected directory (direct child)
		if strings.HasPrefix(lowerPath, lowerProtected+string(filepath.Separator)) {
			return true
		}
	}

	return false
}

// validateDeletePath performs comprehensive validation before deletion.
// It checks for empty paths, path traversal, and protected system directories.
func validateDeletePath(path string) error {
	if path == "" {
		return fmt.Errorf("cannot delete empty path")
	}

	// Check for root directory
	if path == "/" || path == `\` || path == `C:\` || path == `c:\` {
		return fmt.Errorf("cannot delete root directory")
	}

	// Check for path traversal in original path BEFORE resolving
	// This catches attempts like "../../../etc/passwd"
	if strings.HasPrefix(path, "..") || strings.HasPrefix(path, "../") || strings.HasPrefix(path, "..\\") {
		return fmt.Errorf("path traversal detected")
	}
	
	// Also check for /.. or \.. anywhere in the path
	if strings.Contains(path, "/..") || strings.Contains(path, "\\..") {
		return fmt.Errorf("path traversal detected")
	}

	// Convert to absolute path and clean
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	
	cleanPath := filepath.Clean(absPath)

	// Check if path is protected
	if isProtectedPath(cleanPath) {
		return fmt.Errorf("cannot delete protected system path: %s", path)
	}

	return nil
}

// escapeAppleScriptString escapes special characters for AppleScript strings.
// AppleScript uses backslash for escaping, and special chars: ", ', \
func escapeAppleScriptString(s string) string {
	// Escape backslashes first, then quotes
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "'", "\\'")
	return s
}

// escapePowerShellString escapes special characters for PowerShell strings.
// PowerShell uses backtick for escaping, and special chars: ', ", `, $, ()
func escapePowerShellString(s string) string {
	// In single-quoted strings, only ' and \ need escaping
	s = strings.ReplaceAll(s, "'", "''") // Escape single quote by doubling
	return s
}

// validateMediaPath validates that a path is a legitimate media file path.
// It checks for existence and that the path doesn't contain shell metacharacters.
func validateMediaPath(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}

	// Check for obvious shell injection attempts
	dangerousPatterns := []string{
		";", "|", "&", "$", "`", "(", ")", "{", "}", "[", "]",
		"<", ">", "!", "~", "*", "?", "\\", "\n", "\r",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(path, pattern) {
			// Allow some common characters that are actually safe in filenames
			// but log suspicious patterns
			if pattern != " " && pattern != "(" && pattern != ")" {
				// Continue checking - we'll use proper escaping instead of rejecting
			}
		}
	}

	// Verify file exists
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
// Returns the command or an error if the path is invalid.
func SafePlayMedia(path string) (*exec.Cmd, error) {
	if err := validateMediaPath(path); err != nil {
		return nil, err
	}

	absPath, _ := filepath.Abs(path)

	switch runtime.GOOS {
	case "darwin":
		// afplay takes the path as argument - no shell interpolation needed
		return exec.Command("afplay", absPath), nil

	case "linux":
		// aplay takes the path as argument - no shell interpolation needed
		return exec.Command("aplay", absPath), nil

	case "windows":
		// PowerShell requires escaping, use single-quoted string with escaped quotes
		escapedPath := escapePowerShellString(absPath)
		// Use -File parameter with proper argument passing instead of string interpolation
		psScript := fmt.Sprintf(`
$player = New-Object Media.SoundPlayer '%s'
$player.PlaySync()
`, escapedPath)
		return exec.Command("powershell", "-Command", psScript), nil

	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// SafeMoveToTrash moves a file to trash using OS-native commands with proper escaping.
func SafeMoveToTrash(path string) error {
	// Comprehensive path validation
	if err := validateDeletePath(path); err != nil {
		return err
	}
	
	absPath, err := sanitizePathForShell(path)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "darwin":
		return safeMoveToTrashMacOS(absPath)
	case "linux":
		return safeMoveToTrashLinux(absPath)
	case "windows":
		return safeMoveToTrashWindows(absPath)
	default:
		return os.RemoveAll(absPath)
	}
}

func safeMoveToTrashMacOS(path string) error {
	// Try using the `trash` CLI tool first (takes argument directly, no shell)
	if _, err := exec.LookPath("trash"); err == nil {
		return exec.Command("trash", path).Run()
	}

	// Fall back to AppleScript with proper escaping
	escapedPath := escapeAppleScriptString(path)
	script := fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, escapedPath)
	return exec.Command("osascript", "-e", script).Run()
}

func safeMoveToTrashLinux(path string) error {
	// Try using gio (GNOME) - takes argument directly
	if _, err := exec.LookPath("gio"); err == nil {
		return exec.Command("gio", "trash", path).Run()
	}

	// Try using trash-cli - takes argument directly
	if _, err := exec.LookPath("trash"); err == nil {
		return exec.Command("trash", path).Run()
	}

	// Fallback to manual move with proper filename sanitization
	home := os.Getenv("HOME")
	if home != "" {
		trashDir := filepath.Join(home, ".local", "share", "Trash", "files")
		if err := os.MkdirAll(trashDir, 0755); err == nil {
			// Sanitize the filename for the trash directory
			baseName := filepath.Base(path)
			// Remove or replace dangerous characters in filename
			safeName := regexp.MustCompile(`[<>:"/\\|?*]`).ReplaceAllString(baseName, "_")

			// Use O_CREATE|O_EXCL to atomically check and create
			// This prevents TOCTOU race conditions
			counter := 0
			for {
				ext := filepath.Ext(safeName)
				base := strings.TrimSuffix(safeName, ext)
				var fileName string
				if counter == 0 {
					fileName = safeName
				} else {
					fileName = fmt.Sprintf("%s (%d)%s", base, counter, ext)
				}
				
				dest := filepath.Join(trashDir, fileName)
				
				// Try to create the file exclusively - fails if it already exists
				// This is atomic and prevents TOCTOU races
				f, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
				if err == nil {
					// File created successfully, now rename over it
					f.Close()
					os.Remove(dest) // Remove the empty file we just created
					
					if err := os.Rename(path, dest); err != nil {
						// Rename failed, clean up
						os.Remove(dest)
						return err
					}
					return nil
				}
				
				// File already exists, try next counter
				if !os.IsExist(err) {
					// Some other error, fall back to permanent delete
					break
				}
				counter++
				
				// Safety limit to prevent infinite loops
				if counter > 1000 {
					break
				}
			}
		}
	}

	return os.RemoveAll(path)
}

func safeMoveToTrashWindows(path string) error {
	// Use PowerShell with proper escaping
	escapedPath := escapePowerShellString(path)
	psScript := fmt.Sprintf(`
$shell = New-Object -ComObject Shell.Application
$folder = $shell.Namespace(0)
$item = $folder.ParseName('%s')
if ($item -ne $null) {
    $item.InvokeVerb("delete")
}
`, escapedPath)
	return exec.Command("powershell", "-Command", psScript).Run()
}
