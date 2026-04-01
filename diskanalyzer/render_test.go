package diskanalyzer

import (
	"testing"
	"time"
)

func TestRenderCLI(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     1024,
		FileCount:     10,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
		ScannedAt:     time.Now(),
	}

	opts := CLIOptions{
		TopN:   5,
		ByType: true,
	}

	// RenderCLI writes to stdout, just verify it doesn't panic
	RenderCLI(result, opts)
}

func TestRenderCLI_EmptyResult(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     0,
		FileCount:     0,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
		ScannedAt:     time.Now(),
	}

	opts := CLIOptions{
		TopN:   5,
		ByType: false,
	}

	// Should not panic
	RenderCLI(result, opts)
}

func TestRenderCLI_AllOptions(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     1024,
		FileCount:     10,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
		ScannedAt:     time.Now(),
	}

	result.Root.Children = []*DirNode{
		{Path: "/test/dir1", TotalSize: 500, Name: "dir1"},
		{Path: "/test/dir2", TotalSize: 500, Name: "dir2"},
	}

	result.TypeBreakdown[".txt"] = 500
	result.TypeBreakdown[".log"] = 500

	opts := CLIOptions{
		TopN:      5,
		ByType:    true,
		OlderThan: 30,
		MinSize:   100,
		Depth:     2,
	}

	// Should not panic
	RenderCLI(result, opts)
}

func TestRenderTree(t *testing.T) {
	root := &DirNode{
		Path:      "/test",
		TotalSize: 1000,
		Name:      "test",
		Children: []*DirNode{
			{Path: "/test/dir1", TotalSize: 500, Name: "dir1"},
			{Path: "/test/dir2", TotalSize: 500, Name: "dir2"},
		},
	}

	// renderTree writes to stdout, just verify it doesn't panic
	renderTree(root, 2, 2, 1000)
}

func TestRenderTree_LeafNode(t *testing.T) {
	root := &DirNode{
		Path:      "/test/file.txt",
		TotalSize: 100,
		Name:      "file.txt",
		Children:  []*DirNode{},
	}

	// Should not panic
	renderTree(root, 2, 2, 100)
}

func TestRenderTopFiles(t *testing.T) {
	result := &AnalysisResult{
		TotalSize: 1000,
		Root: &DirNode{
			Path:      "/test",
			TotalSize: 1000,
			Children: []*DirNode{
				{Path: "/test/large", TotalSize: 500, Name: "large"},
				{Path: "/test/small", TotalSize: 100, Name: "small"},
			},
		},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
	}

	// Should not panic
	renderTopFiles(result, 5)
}

func TestRenderTopFiles_Empty(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     0,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
	}

	// Should not panic
	renderTopFiles(result, 5)
}

func TestRenderByType(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     1000,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
		ScannedAt:     time.Now(),
	}

	result.TypeBreakdown[".txt"] = 500
	result.TypeBreakdown[".log"] = 300
	result.TypeBreakdown[".jpg"] = 200

	// Should not panic
	renderByType(result)
}

func TestRenderByType_Empty(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     0,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
		ScannedAt:     time.Now(),
	}

	// Should not panic
	renderByType(result)
}

func TestRenderOldFiles(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     1000,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
		ScannedAt:     time.Now(),
	}

	oldTime := time.Now().AddDate(0, 0, -365)
	result.AllFiles = append(result.AllFiles, FileEntry{
		Path:    "/test/old.txt",
		Size:    100,
		ModTime: oldTime,
	})

	// Should not panic
	renderOldFiles(result, 30, 0)
}

func TestRenderOldFiles_None(t *testing.T) {
	result := &AnalysisResult{
		TotalSize:     1000,
		Root:          &DirNode{Path: "/test"},
		AllFiles:      make([]FileEntry, 0),
		TypeBreakdown: make(map[string]int64),
		ScannedAt:     time.Now(),
	}

	// Should not panic
	renderOldFiles(result, 30, 0)
}

func TestMakeBar(t *testing.T) {
	bar := makeBar(50, 100, 20)

	if len(bar) == 0 {
		t.Error("Expected non-empty bar")
	}
}

func TestMakeBar_Full(t *testing.T) {
	bar := makeBar(100, 100, 20)

	if len(bar) < 20 {
		t.Errorf("Expected full bar to be at least 20 chars, got %d", len(bar))
	}
}

func TestMakeBar_Empty(t *testing.T) {
	bar := makeBar(0, 100, 20)

	// Empty bar should still have some representation
	_ = bar
}

func TestGetTerminalWidth(t *testing.T) {
	width := GetTerminalWidth(80)

	if width < 40 {
		t.Errorf("Expected terminal width >= 40, got %d", width)
	}

	if width > 1000 {
		t.Errorf("Expected terminal width <= 1000, got %d", width)
	}
}
