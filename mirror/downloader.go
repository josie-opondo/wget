package mirror

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// expandPath expands shorthand notations to full paths
func expandPath(path string) string {
	// 1. Expand `~` to the home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error finding home directory:", err)
			return ""
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// 2. Expand environment variables like $HOME, $USER, etc.
	path = os.ExpandEnv(path)

	// 3. Convert relative paths (./ or ../) to absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}

	return absPath
}
