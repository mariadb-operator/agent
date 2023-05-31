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
	if _, err := os.Stat(configDir); err != nil {
		return nil, fmt.Errorf("error reading config directory: %v", err)
	}
	if _, err := os.Stat(stateDir); err != nil {
		return nil, fmt.Errorf("error reading state directory: %v", err)
	}
	return &FileManager{
		configDir: configDir,
		stateDir:  stateDir,
	}, nil
}

func (f *FileManager) ReadStateFile(name string) ([]byte, error) {
	return readFile(filepath.Join(f.stateDir, name))
}

func (f *FileManager) WriteStateFile(name string, bytes []byte) error {
	return writeFile(filepath.Join(f.stateDir, name), bytes)
}

func readFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return bytes, nil
}

func writeFile(path string, bytes []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, bytes, info.Mode()); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}
