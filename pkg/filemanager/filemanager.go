package filemanager

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	writeFileMode = fs.FileMode(0777)
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
	return os.ReadFile(filepath.Join(f.stateDir, name))
}

func (f *FileManager) WriteStateFile(name string, bytes []byte) error {
	return os.WriteFile(filepath.Join(f.stateDir, name), bytes, writeFileMode)
}

func (f *FileManager) WriteConfigFile(name string, bytes []byte) error {
	return os.WriteFile(filepath.Join(f.configDir, name), bytes, writeFileMode)
}

func (f *FileManager) DeleteConfigFile(name string) error {
	return os.Remove(filepath.Join(f.configDir, name))
}
