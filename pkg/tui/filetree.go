package tui

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// parentDirEntry implements os.DirEntry for the ".." entry
type parentDirEntry struct{}

func (e parentDirEntry) Name() string               { return ".." }
func (e parentDirEntry) IsDir() bool                { return true }
func (e parentDirEntry) Type() fs.FileMode          { return fs.ModeDir }
func (e parentDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

type FileTreeModel struct {
	path     string
	files    []os.DirEntry
	cursor   int
	selected string
	width    int
	height   int
}

func NewFileTree(startPath string) FileTreeModel {
	m := FileTreeModel{
		path: startPath,
	}
	m.loadFiles()
	return m
}

func (m *FileTreeModel) loadFiles() {
	entries, _ := os.ReadDir(m.path)

	// Prepend ".." if we can go up (ignoring root check for simplicity for now, or check generic root)
	// For robust checking, we could compare filepath.Dir(m.path) with m.path, but usually safe on non-root.
	// Let's just always add it if path is not absolute root (e.g. C:\ or /).

	parent := filepath.Dir(m.path)
	if parent != m.path {
		// Create a slice with space for .. + entries
		newFiles := make([]os.DirEntry, 0, len(entries)+1)
		newFiles = append(newFiles, parentDirEntry{})
		newFiles = append(newFiles, entries...)
		m.files = newFiles
	} else {
		m.files = entries
	}
}

func (m FileTreeModel) Init() tea.Cmd {
	return nil
}

func (m FileTreeModel) Update(msg tea.Msg) (FileTreeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.files) == 0 {
				break
			}
			selected := m.files[m.cursor]

			// Handle ".." specifically or via standard IsDir
			if selected.Name() == ".." {
				m.path = filepath.Dir(m.path)
				m.loadFiles()
				m.cursor = 0
				return m, nil
			}

			newPath := filepath.Join(m.path, selected.Name())
			if selected.IsDir() {
				m.path = newPath
				m.loadFiles()
				m.cursor = 0
			} else {
				// Send generic OpenFileMsg
				m.selected = newPath
				return m, func() tea.Msg {
					return OpenFileMsg{Path: newPath}
				}
			}
		case "backspace", "left", "h": // still support shortcuts
			parent := filepath.Dir(m.path)
			if parent != m.path {
				m.path = parent
				m.loadFiles()
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m FileTreeModel) View() string {
	var s strings.Builder

	// Header
	s.WriteString(StyleDim.Render(fmt.Sprintf("DIR: %s", m.path)) + "\n\n")

	for i, f := range m.files {
		// Custom selection indicator
		if m.cursor == i {
			// Selected: "> filename" in pink
			line := fmt.Sprintf("%s %s", "➜", f.Name())
			if f.IsDir() {
				line += "/"
			}
			s.WriteString(StyleSelected.Render(line) + "\n")
		} else {
			// Unselected: "  filename" in dim
			line := fmt.Sprintf("  %s", f.Name())
			if f.IsDir() {
				line += "/"
			}
			s.WriteString(StyleDim.Render(line) + "\n")
		}
	}

	return StyleFileTree.
		Width(m.width).
		Height(m.height).
		Render(s.String())
}

type OpenFileMsg struct {
	Path string
}
