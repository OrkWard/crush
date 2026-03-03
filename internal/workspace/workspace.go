// Package workspace provides functionality for managing per-project data directories
// in a centralized location, similar to VS Code's workspace storage model.
package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/crush/internal/home"
)

const workspacesDir = "workspaces"

// Manager handles workspace data directory resolution and creation.
type Manager struct {
	baseDir string
}

// NewManager creates a new workspace manager with the default base directory.
// The base directory is typically ~/.local/state/crush/workspaces on Unix
// and %LOCALAPPDATA%/crush/workspaces on Windows.
func NewManager() *Manager {
	return &Manager{
		baseDir: defaultWorkspacesDir(),
	}
}

// NewManagerWithBase creates a new workspace manager with a custom base directory.
func NewManagerWithBase(baseDir string) *Manager {
	return &Manager{
		baseDir: baseDir,
	}
}

// GetDataDir returns the data directory for a given working directory.
// It creates a unique subdirectory based on the working directory path hash.
func (m *Manager) GetDataDir(workingDir string) (string, error) {
	if workingDir == "" {
		return "", fmt.Errorf("working directory cannot be empty")
	}

	// Get absolute path for consistency
	absPath, err := filepath.Abs(workingDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Create a unique hash for this workspace
	hash := hashPath(absPath)

	// Use the last directory name + hash for readability
	base := filepath.Base(absPath)
	if base == "/" || base == "." {
		base = "root"
	}

	workspaceDir := filepath.Join(m.baseDir, fmt.Sprintf("%s-%s", base, hash))

	// Ensure the directory exists
	if err := os.MkdirAll(workspaceDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	return workspaceDir, nil
}

// ListWorkspaces returns all existing workspaces.
func (m *Manager) ListWorkspaces() ([]Workspace, error) {
	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read workspaces directory: %w", err)
	}

	var workspaces []Workspace
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		workspaces = append(workspaces, Workspace{
			Name:    entry.Name(),
			Path:    filepath.Join(m.baseDir, entry.Name()),
			ModTime: info.ModTime(),
		})
	}

	return workspaces, nil
}

// Workspace represents a single workspace directory.
type Workspace struct {
	Name    string
	Path    string
	ModTime interface{}
}

// hashPath creates a SHA256 hash of the path for unique identification.
func hashPath(path string) string {
	hash := sha256.Sum256([]byte(path))
	return hex.EncodeToString(hash[:])[:16]
}

// defaultWorkspacesDir returns the default workspaces directory.
func defaultWorkspacesDir() string {
	// Use the same base as GlobalConfigData but with a workspaces subdirectory
	if crushData := os.Getenv("CRUSH_GLOBAL_DATA"); crushData != "" {
		return filepath.Join(crushData, workspacesDir)
	}

	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "crush", workspacesDir)
	}

	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		return filepath.Join(localAppData, "crush", workspacesDir)
	}

	return filepath.Join(home.Dir(), ".local", "state", "crush", workspacesDir)
}
