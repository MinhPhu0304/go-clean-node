package main

import (
	"flag"
	"fmt"

	"github.com/MinhPhu0304/go-clean-node/cleaner"
)

func main() {
	dryrun := flag.Bool("dryrun", true, "Actually delete it or not")
	path := flag.String("path", "", "Path to your node_modules. Eg: ./foo/bar/node_modules")
	flag.Parse()
	if *path == "" {
		fmt.Println("Missing file path please try again")
		return
	}

	cleaner := cleaner.NewCleaner(*dryrun, *path)
	cleaner.StartClean()
}
