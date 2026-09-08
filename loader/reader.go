package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"github.com/rs/zerolog/log"
)

func readJSONFile[T any](fsys fs.FS, name string) (T, error) {
	var result T
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return result, fmt.Errorf("read %s: %w", name, err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("parse %s: %w", name, err)
	}
	return result, nil
}

func readOptionalJSONFile[T any](fsys fs.FS, name string) (T, error) {
	result, err := readJSONFile[T](fsys, name)
	if errors.Is(err, fs.ErrNotExist) {
		log.Warn().Err(err).Msgf("Optional file %s not found", name)
		var zero T
		return zero, nil
	}
	return result, err
}
