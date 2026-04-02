package cleaner

import (
	"testing"
)

// TestValidateDeletePath_EmptyPath tests that empty paths are rejected
func TestValidateDeletePath_EmptyPath(t *testing.T) {
	err := validateDeletePath("")
	if err == nil {
		t.Error("validateDeletePath(\"\") should return error")
	}
	if err != nil && err.Error() != "cannot delete empty path" {
		t.Errorf("Expected 'cannot delete empty path', got: %v", err)
	}
}

// TestValidateDeletePath_RootDirectory tests that root directories are rejected
func TestValidateDeletePath_RootDirectory(t *testing.T) {
	rootPaths := []string{"/", `\`, `C:\`, `c:\`}
	
	for _, path := range rootPaths {
		err := validateDeletePath(path)
		if err == nil {
			t.Errorf("validateDeletePath(%q) should return error", path)
		}
	}
}

// TestValidateDeletePath_ProtectedSystemPaths tests that system paths are protected
func TestValidateDeletePath_ProtectedSystemPaths(t *testing.T) {
	protectedPaths := []string{
		// Unix/Linux
		"/etc",
		"/etc/passwd",
		"/bin",
		"/bin/bash",
		"/usr",
		"/usr/bin",
		"/lib",
		"/lib64",
		"/boot",
		"/dev",
		"/proc",
		"/sys",
		// macOS
		"/System",
		"/System/Library",
		"/Library",
		"/Library/Frameworks",
		"/Applications",
		"/Applications/Safari.app",
	}

	for _, path := range protectedPaths {
		err := validateDeletePath(path)
		if err == nil {
			t.Errorf("validateDeletePath(%q) should return error for protected path", path)
		}
	}
}

// TestValidateDeletePath_AllowedSubPaths tests that allowed subpaths are NOT protected
func TestValidateDeletePath_AllowedSubPaths(t *testing.T) {
	// These should NOT be rejected (they're explicitly allowed)
	allowedPaths := []string{
		"/var/tmp",
		"/var/tmp/test.txt",
		"/var/folders",
		"/var/folders/abc123",
		"/var/log",
		"/var/log/syslog",
		"/usr/local",
		"/usr/local/bin",
	}

	for _, path := range allowedPaths {
		err := validateDeletePath(path)
		// Should not be rejected as protected (may fail existence check)
		if err != nil && err.Error() == "cannot delete protected system path: "+path {
			t.Errorf("validateDeletePath(%q) should not reject allowed path", path)
		}
	}
}

// TestValidateDeletePath_PathTraversal tests path traversal detection
func TestValidateDeletePath_PathTraversal(t *testing.T) {
	traversalPaths := []string{
		"../../../etc/passwd",
		"/home/user/../../../etc/passwd",
		"foo/../../../etc/passwd",
	}

	for _, path := range traversalPaths {
		err := validateDeletePath(path)
		if err == nil {
			t.Errorf("validateDeletePath(%q) should detect path traversal", path)
		}
	}
}

// TestIsProtectedPath tests the isProtectedPath helper function
func TestIsProtectedPath(t *testing.T) {
	tests := []struct {
		path      string
		protected bool
	}{
		// Protected paths
		{"/etc", true},
		{"/etc/passwd", true},
		{"/usr/bin", true},
		{"/System/Library", true},
		// Allowed paths (NOT protected)
		{"/var/tmp", false},
		{"/var/tmp/test", false},
		{"/var/folders", false},
		{"/var/log", false},
		{"/usr/local", false},
		// Normal paths
		{"/home/user", false},
		{"/tmp", false},
		{"/Users/user", false},
	}

	for _, tt := range tests {
		result := isProtectedPath(tt.path)
		if result != tt.protected {
			t.Errorf("isProtectedPath(%q) = %v, want %v", tt.path, result, tt.protected)
		}
	}
}

// TestIsProtectedPath_CaseInsensitive tests case-insensitive matching
func TestIsProtectedPath_CaseInsensitive(t *testing.T) {
	tests := []struct {
		path      string
		protected bool
	}{
		{"/ETC", true},
		{"/Etc", true},
		{"/SYSTEM", true},
		{"/system", true},
	}

	for _, tt := range tests {
		result := isProtectedPath(tt.path)
		if result != tt.protected {
			t.Errorf("isProtectedPath(%q) = %v, want %v", tt.path, result, tt.protected)
		}
	}
}

// TestSafeMoveToTrash_ProtectedPath tests that SafeMoveToTrash rejects protected paths
func TestSafeMoveToTrash_ProtectedPath(t *testing.T) {
	protectedPaths := []string{
		"/etc",
		"/bin",
		"/usr",
		"/System",
	}

	for _, path := range protectedPaths {
		err := SafeMoveToTrash(path)
		if err == nil {
			t.Errorf("SafeMoveToTrash(%q) should return error for protected path", path)
		}
	}
}

// TestDeleteEntry_ProtectedPath tests that deleteEntry rejects protected paths
func TestDeleteEntry_ProtectedPath(t *testing.T) {
	entry := EntryInfo{
		Path: "/etc/passwd",
		Size: 100,
	}

	_, _, _, err := deleteEntry(entry, true)
	if err == nil {
		t.Error("deleteEntry should reject protected paths")
	}
}

// TestDeleteEntry_PathTraversal tests that deleteEntry detects path traversal
func TestDeleteEntry_PathTraversal(t *testing.T) {
	entry := EntryInfo{
		Path: "../../../etc/passwd",
		Size: 100,
	}

	_, _, _, err := deleteEntry(entry, true)
	if err == nil {
		t.Error("deleteEntry should detect path traversal")
	}
}
