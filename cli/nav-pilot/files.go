package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil))[:16], nil
}

// normalizeMarkdown normalizes markdown content for comparison:
//   - CRLF → LF
//   - Trim trailing whitespace per line
//   - Collapse consecutive blank lines to a single blank line
func normalizeMarkdown(data []byte) []byte {
	// CRLF → LF
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

	lines := bytes.Split(data, []byte("\n"))
	var out [][]byte
	prevBlank := false
	for _, line := range lines {
		trimmed := bytes.TrimRight(line, " \t")
		blank := len(trimmed) == 0
		if blank && prevBlank {
			continue
		}
		out = append(out, trimmed)
		prevBlank = blank
	}
	return bytes.Join(out, []byte("\n"))
}

// normalizedFileHash hashes a file after normalizing markdown content.
// For non-.md files, falls back to raw fileHash.
func normalizedFileHash(path string) (string, error) {
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return fileHash(path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	normalized := normalizeMarkdown(data)
	h := sha256.New()
	h.Write(normalized)
	return hex.EncodeToString(h.Sum(nil))[:16], nil
}

// dirHash hashes all files in a directory recursively.
// Markdown files (.md) are normalized before hashing for formatting tolerance.
func dirHash(dir string) (string, error) {
	h := sha256.New()
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		h.Write([]byte(rel))

		if strings.HasSuffix(strings.ToLower(rel), ".md") {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			h.Write(normalizeMarkdown(data))
		} else {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(h, f); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil))[:16], nil
}

// copyFile copies a single file atomically, creating parent directories.
// Refuses to overwrite symlinks to prevent writing outside the repo.
func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	// B2: Refuse to write through symlinks (file or parent directory)
	if err := checkSymlink(dst); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// I4: Atomic write via temp file + rename
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".nav-pilot-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(tmp, in); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	return os.Rename(tmpPath, dst)
}

// checkSymlink detects symlinks in the path chain.
// Checks both the target file and its parent directory.
func checkSymlink(path string) error {
	// Check the file itself if it exists
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to overwrite symlink: %s", path)
	}
	// Check the parent directory
	parent := filepath.Dir(path)
	resolved, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return nil // parent doesn't exist yet, will be created by MkdirAll
	}
	if resolved != parent {
		absParent, _ := filepath.Abs(parent)
		if resolved != absParent {
			return fmt.Errorf("refusing to write through symlinked directory: %s -> %s", parent, resolved)
		}
	}
	return nil
}

// copyDir copies a directory recursively, creating it fresh (removes stale files).
func copyDir(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func countDirFiles(dir string) int {
	count := 0
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			count++
		}
		return nil
	})
	return count
}

// ─── Conflict detection ─────────────────────────────────────────────────────

type conflict struct {
	Path    string
	Current string // hash of existing file
	New     string // hash of source file
}

func checkConflict(targetPath, sourcePath string, isDir bool) (*conflict, error) {
	if isDir {
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			return nil, nil
		}
		currentHash, err := dirHash(targetPath)
		if err != nil {
			return nil, err
		}
		newHash, err := dirHash(sourcePath)
		if err != nil {
			return nil, err
		}
		if currentHash == newHash {
			return nil, nil
		}
		return &conflict{Path: targetPath, Current: currentHash, New: newHash}, nil
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return nil, nil
	}
	currentHash, err := fileHash(targetPath)
	if err != nil {
		return nil, fmt.Errorf("hashing %s: %w", targetPath, err)
	}
	newHash, err := fileHash(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("hashing %s: %w", sourcePath, err)
	}
	if currentHash == newHash {
		return nil, nil
	}
	return &conflict{Path: targetPath, Current: currentHash, New: newHash}, nil
}
