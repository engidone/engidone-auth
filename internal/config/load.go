package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"

	"engidoneauth/log"
)

func LoadFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Error("failed to unmarshal yaml: %v", err)
		return nil, err
	}

	return &config, nil
}

func NewConfigPaths(levels ...string) Paths {

	_, b, _, _ := runtime.Caller(0)

	args := append([]string{filepath.Dir(b)}, levels...)
	rootPath := filepath.Join(args...)
	return Paths{
		Config: filepath.Join(rootPath, "cmd", "config"),
		Root:   rootPath,
	}
}
