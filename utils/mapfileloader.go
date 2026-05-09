package utils

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// Find JSON files in the given `path`
// Do not enter subdirectories as json meta data is expected to live
// at `path`
func GatherMapMetaDataFiles(path string) []string {
	LogDebug("Loading Files from %s", path)
	var result []string

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			LogError("Skipping directory %+v due to error: %+v", d, err)
			return fs.SkipDir
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			result = append(result, p)
		}

		return nil
	})

	if err != nil {
		LogError("Error walking path %s: %+v", path, err)
	}

	return result
}
