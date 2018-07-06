package main

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"strings"
)

func glob(dir string, ext string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func emptyDir(directory string) error {
	err := os.RemoveAll(directory)
	if err != nil {
		return err
	}
	err = os.MkdirAll(directory, 0700)
	return err
}

func searchAndReplaceInFile(file string, old string, new string) error {
	read, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	newContents := strings.Replace(string(read), old, new, -1)
	err = ioutil.WriteFile(file, []byte(newContents), 0644)
	return err
}

func trimLeftComment(s string) string {
	return strings.TrimLeft(s, "/ \t")
}