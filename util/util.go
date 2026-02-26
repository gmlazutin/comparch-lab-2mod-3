package util

import (
	"os"
	"path/filepath"
)

func ListFilesWithExts(dir string, exts []string) ([]string, error) {
	var files []string

	extMap := make(map[string]bool)
	for _, e := range exts {
		extMap[e] = true
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if extMap[filepath.Ext(entry.Name())] {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return files, nil
}
