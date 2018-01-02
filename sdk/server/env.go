package server

import (
	"path/filepath"
	"os"
	"fmt"
	"runtime"
)

type Env struct{
	RootDir       string
	StaticUrl     string
	StaticRoot	  string
}

func NewEnv() *Env {
	rootDir := getRootDir()
	if rootDir == "" {
		return nil
	}

	env := &Env{}
	env.RootDir = rootDir
	env.StaticUrl = "/static/"
	env.StaticRoot = filepath.Join(rootDir, "static")
	return env
}

//getRootDir returns the current path to executable, e,g: wegigo
func getRootDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error get current execute path: %s", err)
		return ""
	}

	return dir
}

func getCurrentDir() (string, bool) {
	//curFilePath is: /path/to/go-wegigo/env.go
	_, curFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", false
	}

	return filepath.Dir(curFilePath), true
}

func getSubDir(rootDir string, dirname string) string {
	subdir := filepath.Join(rootDir, dirname)
	stat, err := os.Stat(subdir)
	if err != nil {
		return ""
	}
	if !stat.IsDir() {
		return ""
	}

	return subdir
}
