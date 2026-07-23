// Package motionsource resolves one-sheet preview sources and complete motion
// source families for the asset importer and release validator.
package motionsource

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveSetPaths accepts either one standalone/set00 preview source or a
// complete set00 through setNN family. A partially populated family is an
// error so import and release validation cannot silently duplicate set00.
func ResolveSetPaths(set00Path string, setCount int) ([]string, error) {
	if set00Path == "" {
		return nil, fmt.Errorf("motion source path is empty")
	}
	if setCount < 1 {
		return nil, fmt.Errorf("motion source set count must be positive, got %d", setCount)
	}
	if !strings.Contains(set00Path, "set00") {
		if _, err := os.Stat(set00Path); err != nil {
			return nil, err
		}
		return []string{set00Path}, nil
	}

	paths := make([]string, 0, setCount)
	missing := make([]string, 0, setCount)
	for set := 0; set < setCount; set++ {
		path := strings.Replace(set00Path, "set00", fmt.Sprintf("set%02d", set), 1)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				missing = append(missing, filepath.Base(path))
				continue
			}
			return nil, err
		}
		paths = append(paths, path)
	}

	switch len(paths) {
	case setCount:
		return paths, nil
	case 1:
		if paths[0] == set00Path {
			return []string{set00Path}, nil
		}
	}
	lastSet := fmt.Sprintf("set%02d", setCount-1)
	return nil, fmt.Errorf(
		"incomplete motion source family for %s: found %d of %d sheets; provide only set00 for preview fallback or all set00-%s (missing: %s)",
		set00Path,
		len(paths),
		setCount,
		lastSet,
		strings.Join(missing, ", "),
	)
}
