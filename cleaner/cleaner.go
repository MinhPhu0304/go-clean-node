package cleaner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/bytefmt"
)

type Cleaner struct {
	dryRun    bool
	path      string
	verbose   bool
	sizeCount int64
}

var junkFiles = []string{
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
	"makefile",
	".nycrc",
	".markdown",
	".jshintrc",
	".codeclimate.yml",
	"editorconfig",
	".github/",
	".coveralls.yml",
	".nyc_output",
	".circleci",
	".vscode",
}

func NewCleaner(dryRun bool, path string, verbose bool) Cleaner {
	return Cleaner{
		dryRun:  dryRun,
		path:    path,
		verbose: verbose,
	}
}

func (c *Cleaner) StartClean() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Failed to clean node_modules")
		}
	}()
	err := filepath.WalkDir(c.path, c.cleanUp)
	if err != nil {
		fmt.Printf("%v", fmt.Errorf("Error cleaning up node_modules: ", err))
	} else if c.dryRun {
		fmt.Println("\nDry run complete\nPass --dryrun=false to delete these files")
	} else {
		fmt.Printf("Total delete size: %sB \n", bytefmt.ByteSize(uint64(c.sizeCount)))
	}
}

func (c *Cleaner) cleanUp(path string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		return nil
	}
	for _, junkFile := range junkFiles {
		if strings.Contains(path, junkFile) {
			info, err := d.Info()
			if err != nil {
				// there's a race condition somewhere in the walkdir so the file can appear twice
				// ¯\_(ツ)_/¯
				return nil
			}
			c.sizeCount = c.sizeCount + info.Size()
			if c.dryRun {
				fmt.Printf("Clean file %s\n", path)
			} else {
				c.deleteFile(path)
			}
		}
	}
	return nil
}

func (c *Cleaner) deleteFile(path string) {
	if c.verbose {
		fmt.Printf("Deletef file with path %s\n", path)
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}
