# Test Plan for DupClean

## Current Coverage Status

| Package | Coverage | Status |
|---------|----------|--------|
| scanner | 55.0% | ✅ Good |
| diskanalyzer | 37.5% | ⚠️ Needs improvement |
| ui | 33.3% | ⚠️ Needs improvement |
| gui | 4.0% | ❌ Critical - needs tests |
| cleaner | 0.0% | ❌ Critical - needs tests |
| internal/fsutil | 0.0% | ❌ Needs tests |

## Priority 1: Cleaner Package (Critical)

The cleaner package is used by the GUI cache cleaner feature but has no tests.

### Tests to Add:
- [ ] `cleaner/scanner_test.go`
  - Test Scan() with various options
  - Test progress callback
  - Test with non-existent paths
  - Test with empty directories
  
- [ ] `cleaner/deleter_test.go`
  - Test Delete() with dry-run mode
  - Test Delete() with permanent deletion
  - Test moveToTrash() on different OS
  - Test isFileInUse() detection
  - Test error handling for permission denied

- [ ] `cleaner/targets_test.go`
  - Test GetSystemTargets() on different OS
  - Test GetBrowserTargets() on different OS
  - Test GetDeveloperTargets() on different OS
  - Test GetLogsTargets() on different OS
  - Test FilterTargets() with various filters

## Priority 2: GUI Package (Critical)

The GUI has complex logic but minimal test coverage.

### Tests to Add:
- [ ] `gui/cache_cleaner_test.go`
  - Test CacheCleanerWidget creation
  - Test startCacheScan() with mock data
  - Test displayCacheResults() rendering
  - Test updateCacheTotal() calculations
  - Test cleanPath() with various patterns
  - Test isProtectedPath() for different paths
  
- [ ] `gui/sidebar_test.go`
  - Test Sidebar() creation
  - Test CreateSidebar() navigation
  - Test item selection highlighting
  
- [ ] `gui/duplicate_finder_test.go`
  - Test DuplicateFinderWidget() creation
  - Test DuplicateResultsWidget() navigation
  - Test ShowDuplicateResults() logic

## Priority 3: Internal/fsutil

- [ ] `internal/fsutil/measure_test.go`
  - Test MeasureDir() with empty directory
  - Test MeasureDir() with files
  - Test MeasureDir() with patterns
  - Test MeasureDir() with minAge filter
  - Test computeDirSize() accuracy

## Priority 4: Improve Existing Coverage

### Scanner Package (55% → 80%)
- [ ] Test ByteScanner.Scan()
- [ ] Test PhotoScanner.Scan() (if exists)
- [ ] Test hashFilePartial() edge cases
- [ ] Test hashFileFull() with large files
- [ ] Test filesIdentical() with different scenarios

### DiskAnalyzer Package (37.5% → 80%)
- [ ] Test error handling in Walk()
- [ ] Test with symlinks
- [ ] Test ExportJSON() output format
- [ ] Test RenderCLI() output formatting

### UI Package (33.3% → 80%)
- [ ] Test CLI input parsing
- [ ] Test invalid user input handling
- [ ] Test moveToTrash() integration

## Test Infrastructure Improvements

- [ ] Add test helpers/mocks for:
  - File system operations
  - OS-specific functions (trash, permissions)
  - Time-based operations
  
- [ ] Add integration tests:
  - Full scan → review → delete workflow
  - Cache cleaner workflow
  - Disk analyzer workflow

- [ ] Add benchmark tests:
  - Benchmark large directory scans
  - Benchmark hashing performance
  - Benchmark UI rendering with many items
