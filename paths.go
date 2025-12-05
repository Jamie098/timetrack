package main

import (
	"os"
	"path/filepath"
)

func getDataDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".timetrack")
	os.MkdirAll(dir, 0755)
	return dir
}

func getDataPath() string {
	return filepath.Join(getDataDir(), "data.json")
}

func getConfigPath() string {
	return filepath.Join(getDataDir(), "config.json")
}

func getPidPath() string {
	return filepath.Join(getDataDir(), "daemon.pid")
}
