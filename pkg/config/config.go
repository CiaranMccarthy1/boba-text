package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Colors struct {
	Text      string `toml:"text"`
	SubText   string `toml:"subtext"`
	Primary   string `toml:"primary"`
	Secondary string `toml:"secondary"`
	Accent    string `toml:"accent"`
	Success   string `toml:"success"`
	Warning   string `toml:"warning"`
	Error     string `toml:"error"`
	Dark      string `toml:"dark"`
}

type Keys struct {
	ToggleTree string `toml:"toggle_tree"`
	FocusTree  string `toml:"focus_tree"`
	FocusAgent string `toml:"focus_agent"`
	CycleFocus string `toml:"cycle_focus"`
	Quit       string `toml:"quit"`
}

type AI struct {
	Name  string `toml:"name"`
	Model string `toml:"model"`
}

type Commands struct {
	Save []string `toml:"save"`
	Quit []string `toml:"quit"`
}

type Config struct {
	Colors   Colors   `toml:"colors"`
	Keys     Keys     `toml:"keys"`
	AI       AI       `toml:"ai"`
	Commands Commands `toml:"commands"`
}

func DefaultConfig() Config {
	return Config{
		Colors: Colors{
			Text:      "#FAFAFA",
			SubText:   "#7D7D7D",
			Primary:   "#F25D94", // Neon Pink
			Secondary: "#A550DF", // Purple
			Accent:    "#61AFEF", // Blue
			Success:   "#98C379", // Green
			Warning:   "#E5C07B", // Yellow
			Error:     "#E06C75", // Red
			Dark:      "#1E1E1E",
		},
		Keys: Keys{
			ToggleTree: "ctrl+b",
			FocusTree:  "ctrl+e",
			FocusAgent: "ctrl+a",
			CycleFocus: "tab",
			Quit:       "ctrl+c",
		},
		AI: AI{
			Name:  "Agent",
			Model: "default",
		},
		Commands: Commands{
			Save: []string{"w", "s", "save", "write"},
			Quit: []string{"q", "quit", "exit"},
		},
	}
}

func Load() Config {
	// Try local file first, then home dir
	paths := []string{"config.toml"}

	home, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths, filepath.Join(home, ".boba-config.toml"))
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			var cfg Config
			// Start with defaults to fill missing fields
			defaults := DefaultConfig()
			cfg = defaults

			// Unmarshal over defaults (this is a simple way, though slightly imperfect for deep structs, works here)
			if _, err := toml.DecodeFile(path, &cfg); err == nil {
				return cfg
			}
		}
	}

	return DefaultConfig()
}
