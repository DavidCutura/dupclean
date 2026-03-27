package cleaner

import (
	"time"
)

// Risk indicates how safe it is to delete a target automatically.
type Risk uint8

const (
	RiskSafe Risk = iota // always safe: OS temp dirs, browser caches
	RiskLow              // safe for most users: app caches, log files
	RiskModerate         // may cause slowdowns on next launch: large app caches
	RiskHigh             // use caution: system-level files, package caches
)

// CleanTarget is a single cleanable location.
type CleanTarget struct {
	ID          string   // unique, stable key e.g. "macos-user-cache"
	Category    string   // display grouping e.g. "System", "Browser", "Developer"
	Label       string   // human label e.g. "User application cache"
	Description string   // one sentence explaining what this is
	Paths       []string // resolved absolute paths to scan (may be empty if not found)
	Patterns    []string // glob patterns within each path e.g. "*.log", "*"
	Risk        Risk
	OS          string // "darwin", "linux", "windows", or "" for all

	// populated after Scan():
	TotalSize   int64
	FileCount   int
	ScannedAt   time.Time
	Entries     []EntryInfo // individual files/dirs found
	Selected    bool        // for CLI/GUI selection state
}

// EntryInfo is one file or directory within a target.
type EntryInfo struct {
	Path    string
	Size    int64     // for dirs: recursive sum
	ModTime time.Time
	IsDir   bool
}

// ScanResult is the output of a full scan run.
type ScanResult struct {
	Targets   []*CleanTarget
	TotalSize int64
	ScannedAt time.Time
	Errors    []ScanError
}

// ScanError is a non-fatal error encountered while scanning a single target.
type ScanError struct {
	TargetID string
	Path     string
	Err      error
}

// Registry returns all targets appropriate for the current OS,
// with Paths resolved using os.UserHomeDir and os.Getenv.
func Registry() []*CleanTarget {
	var all []*CleanTarget

	// Add targets from all categories
	all = append(all, GetSystemTargets()...)
	all = append(all, GetBrowserTargets()...)
	all = append(all, GetDeveloperTargets()...)
	all = append(all, GetLogsTargets()...)

	return all
}
