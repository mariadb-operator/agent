package filemanager

import (
	"fmt"
	"os"
)

type FileManager struct {
	configDir string
	stateDir  string
}

func NewFileManager(configDir, stateDir string) (*FileManager, error) {
	if err := mustExist(configDir); err != nil {
		return nil, err
	}
	if err := mustExist(stateDir); err != nil {
		return nil, err
	}
	return &FileManager{
		configDir: configDir,
		stateDir:  stateDir,
	}, nil
}

func mustExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("'%s' does not exist", path)
	}
	return err
}
