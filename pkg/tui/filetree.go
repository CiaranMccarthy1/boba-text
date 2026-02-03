package tui

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// parentDirEntry implements os.DirEntry for the ".." parent directory entry.
type parentDirEntry struct{}

// Name returns ".." to represent the parent directory.
func (e parentDirEntry) Name() string { return ".." }

// IsDir returns true since parent directory entries are always directories.
func (e parentDirEntry) IsDir() bool { return true }

// Type returns fs.ModeDir indicating a directory type.
func (e parentDirEntry) Type() fs.FileMode { return fs.ModeDir }

// Info returns nil as parent directory entries don't need file info.
func (e parentDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

type FileTreeModel struct {
	path     string
	files    []os.DirEntry
	cursor   int
	selected string
	width    int
	height   int
}

// NewFileTree creates a new file tree model starting at the given path.
func NewFileTree(startPath string) FileTreeModel {
	m := FileTreeModel{
		path: startPath,
	}
	m.loadFiles()
	return m
}

// loadFiles loads directory entries from the current path and adds parent entry if not at root.
func (m *FileTreeModel) loadFiles() {
	entries, _ := os.ReadDir(m.path)

	parent := filepath.Dir(m.path)
	if parent != m.path {
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

// Update handles messages and updates the file tree state.
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
				m.selected = newPath
				return m, func() tea.Msg {
					return OpenFileMsg{Path: newPath}
				}
			}
		case "backspace", "left", "h":
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

// View renders the file tree as a string.
func (m FileTreeModel) View() string {
	var s strings.Builder

	s.WriteString(StyleDim.Render(fmt.Sprintf("DIR: %s", m.path)) + "\n\n")

	for i, f := range m.files {
		if m.cursor == i {
			line := fmt.Sprintf("%s %s", "➜", f.Name())
			if f.IsDir() {
				line += "/"
			}
			s.WriteString(StyleSelected.Render(line) + "\n")
		} else {
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
