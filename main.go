package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/bytefmt"
)

var sizeCount int64 = 0

func main() {
	filepath.WalkDir("./node_modules", cleanUp)
	fmt.Printf("Total delete size: %sB \n", bytefmt.ByteSize(uint64(sizeCount)))
}

var deleteFiles = []string{
	".md",
	"LICENSE",
	"license",
	".txt",
	".test.",
	".spec.",
	"package.json",
	"Jenkins",
	".babelrc",
	".npmrc",
	".eslintrc",
	".eslintignore",
	"Dockerfile",
	".travis.yml",
	".zuul.yml",
	".gitmodules",
	".npmignore",
	"Makefile",
	".nycrc",
	".markdown",
	".jshintrc",
	".codeclimate.yml",
	"editorconfig",
	".github/",
	".coveralls.yml",
}

func cleanUp(path string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		return nil
	}
	for _, junkFile := range deleteFiles {
		if strings.Contains(path, junkFile) {
			info, err := d.Info()
			if err != nil {
				// there's a race condition somewhere in the walkdir so the file can appear twice
				// ¯\_(ツ)_/¯
				return nil
			}
			sizeCount = sizeCount + info.Size()
			deleteFile(path)
		}
	}
	return nil
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}
