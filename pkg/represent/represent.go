// Copyright 2013 Casey Marshall <casey.marshall@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package represent

import (
	"go/build"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const REPRESENT_PKG = "github.com/cmars/represent"
const PRESENT_PKG = "golang.org/x/tools/present"

// Represent provides functions to publish a directory tree of files
// in the Present format.
type Represent struct {
	srcDir        string
	publishDir    string
	presentPkgDir string
}

// NewRepresent builds a new, initialized Represent struct
// for processing the given base directory. If srcDir is
// the empty string, base directory is assumed to be the
// current working directory.
func NewRepresent(srcDir, publishDir string) (*Represent, error) {
	// Discern where the represent package source is, so we can
	// locate the static files.
	p, err := build.Default.Import(REPRESENT_PKG, "", build.FindOnly)
	if err != nil {
		return nil, err
	}
	// Resolve default source directory
	if srcDir == "" {
		srcDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	// Resolve default publish directory
	if publishDir == "" {
		publishDir = filepath.Join(srcDir, "publish")
	}
	r := &Represent{srcDir: srcDir,
		publishDir:    publishDir,
		presentPkgDir: p.Dir}
	err = r.requirePublishDir()
	return r, err
}

// requirePublishDir creates the publish directory
// if it does not exist.
func (r *Represent) requirePublishDir() (err error) {
	err = os.MkdirAll(r.srcDir, 0755)
	if err != nil {
		return
	}
	return os.MkdirAll(r.publishDir, 0755)
}

// Publish sets up all referenced static assets, and compiles
// all Present format files in the source directory tree to
// the publish directory.
func (r *Represent) Publish() (err error) {
	err = r.updateAssets()
	if err != nil {
		return
	}
	return r.compileAll()
}

// updateAssets copies static files (Javascript, CSS, etc)
// that are referenced by the Present templates.
func (r *Represent) updateAssets() error {
	staticDir := filepath.Join(r.presentPkgDir, "static")
	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Derive relative path to static file
		relPath := strings.Replace(path, r.presentPkgDir, "", 1)
		targetPath := filepath.Join(append([]string{r.publishDir}, filepath.SplitList(relPath)...)...)
		// Create destination parent dir if not exists
		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return err
		}
		// Copy file to destination
		destFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

// compileAll publishes all the .slide and .article files found
// in the base directory into the publish directory.
func (r *Represent) compileAll() error {
	return filepath.Walk(r.srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Clean(path) == filepath.Clean(r.publishDir) {
			// Skip the publish dir
			return filepath.SkipDir
		}
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}
		// Derive relative path to static file
		relPath := strings.Replace(path, r.srcDir, "", 1)
		targetPath := filepath.Join(append([]string{r.publishDir}, filepath.SplitList(relPath)...)...)
		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return err
		}
		if isDoc(path) {
			// Render Present file to static HTML
			targetPath = regexp.MustCompile(`\.[^.]+$`).ReplaceAllString(targetPath, ".html")
			destFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer destFile.Close()
			return renderDoc(destFile, r.presentPkgDir, path)
		}
		// Copy file into publish directory unchanged
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		destFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, srcFile)
		return err
	})
}
