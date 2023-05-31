package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileManager struct {
	configDir string
	stateDir  string
}

func NewFileManager(configDir, stateDir string) (*FileManager, error) {
	if err := fileMustExist(configDir); err != nil {
		return nil, fmt.Errorf("config directory does not exist: %v", err)
	}
	if err := fileMustExist(stateDir); err != nil {
		return nil, fmt.Errorf("state directory does not exist: %v", err)
	}
	return &FileManager{
		configDir: configDir,
		stateDir:  stateDir,
	}, nil
}

func (f *FileManager) ReadStateFile(name string) ([]byte, error) {
	return readFile(filepath.Join(f.stateDir, name))
}

func readFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil, err
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return bytes, nil
}

func fileMustExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("'%s' does not exist", path)
	}
	return err
}
