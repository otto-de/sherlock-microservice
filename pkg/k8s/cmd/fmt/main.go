package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"sigs.k8s.io/yaml/kyaml"
)

var (
	chdir = flag.String("C", "", "directory to search for YAML files")
)

func main() {
	flag.Parse()

	if *chdir == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		chdir = &wd
	}

	formatting := sync.WaitGroup{}

	err := formatKYAMLFiles(*chdir, &formatting, formatKYAML)
	if err != nil {
		panic(err)
	}

	formatting.Wait()
}

// formatKYAMLFiles recursively searches for KYAML files in the given directory
func formatKYAMLFiles(root string, wg *sync.WaitGroup, formatFunc func(string) error) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".kyaml" {
			wg.Go(func() {
				err := formatFunc(path)
				if err != nil {
					panic(err)
				}
			})
		}

		return nil
	})

	return err
}

// formatKYAML reads a KYAML file into memory, formats it, and writes it back
func formatKYAML(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	formatted, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer formatted.Close()

	r := bytes.NewReader(data)

	enc := kyaml.Encoder{}
	err = enc.FromYAML(r, formatted)
	if err != nil {
		return fmt.Errorf("failed to reencode KYAML: %w", err)
	}

	return nil
}
