package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// isInteractive returns true when stdin is a terminal (not piped).
func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// isGitRepo returns true if dir contains a .git directory.
func isGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// cmdInteractive runs an interactive collection picker and installer.
func cmdInteractive() error {
	targetDir, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	// If already installed, offer sync instead
	state, _ := readState(targetDir)
	if state != nil {
		fmt.Printf("%s is already installed %s\n\n",
			bold(state.Collection), dim(fmt.Sprintf("(v%s, %s)", state.Version, state.SourceSHA)))
		fmt.Printf("Run %s to check for updates.\n", bold("nav-pilot sync"))
		return nil
	}

	fmt.Println(bold("nav-pilot") + dim(" — Nav's Copilot toolkit"))
	fmt.Println()
	fmt.Println(dim("Resolving source..."))

	src, err := resolveSource("", "")
	if err != nil {
		return err
	}
	defer src.Cleanup()

	names, err := listCollectionDirs(src.Dir)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("no collections found")
	}

	// Display collections
	type collectionInfo struct {
		name  string
		desc  string
		total int
	}
	var collections []collectionInfo
	for _, name := range names {
		m, err := loadManifest(src.Dir, name)
		if err != nil {
			continue
		}
		total := len(m.Agents) + len(m.Skills) + len(m.Instructions) + len(m.Prompts)
		collections = append(collections, collectionInfo{name: name, desc: m.Description, total: total})
	}

	if len(collections) == 0 {
		return fmt.Errorf("no valid collections found")
	}

	fmt.Println()
	fmt.Println(bold("Available collections:"))
	fmt.Println()
	for i, c := range collections {
		fmt.Printf("  %s  %-20s %s %s\n",
			bold(fmt.Sprintf("%d.", i+1)),
			c.name,
			c.desc,
			dim(fmt.Sprintf("(%d items)", c.total)))
	}
	fmt.Println()

	// Prompt for selection
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Select collection [1-%d]: ", len(collections))
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(collections) {
		return fmt.Errorf("invalid selection: %q", input)
	}
	selected := collections[choice-1]

	// Show preview
	fmt.Println()
	m, err := loadManifest(src.Dir, selected.name)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s — %s\n", dim("→"), bold(selected.name), m.Description)
	parts := []string{}
	if len(m.Agents) > 0 {
		parts = append(parts, fmt.Sprintf("%d agents", len(m.Agents)))
	}
	if len(m.Skills) > 0 {
		parts = append(parts, fmt.Sprintf("%d skills", len(m.Skills)))
	}
	if len(m.Instructions) > 0 {
		parts = append(parts, fmt.Sprintf("%d instructions", len(m.Instructions)))
	}
	if len(m.Prompts) > 0 {
		parts = append(parts, fmt.Sprintf("%d prompts", len(m.Prompts)))
	}
	fmt.Printf("  %s\n", dim(strings.Join(parts, ", ")))
	fmt.Println()

	// Confirm
	fmt.Printf("Install %s? [Y/n]: ", bold(selected.name))
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "" && confirm != "y" && confirm != "yes" {
		fmt.Println(dim("Cancelled."))
		return nil
	}

	// Install
	fmt.Println()
	return cmdInstall(selected.name, targetDir, "", false, false)
}
