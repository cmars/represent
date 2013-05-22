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
const PRESENT_PKG = "code.google.com/p/go.talks/present"

type Represent struct {
	BaseDir       string
	publishDir    string
	presentPkgDir string
}

func NewRepresent(baseDir string) (*Represent, error) {
	p, err := build.Default.Import(REPRESENT_PKG, "", build.FindOnly)
	if err != nil {
		return nil, err
	}
	r := &Represent{BaseDir: baseDir,
		publishDir:    filepath.Join(baseDir, "publish"),
		presentPkgDir: p.Dir}
	return r, nil
}

func (r *Represent) RequirePublishDir() (err error) {
	err = os.MkdirAll(r.BaseDir, 0755)
	if err != nil {
		return
	}
	err = os.RemoveAll(r.publishDir)
	if err != nil {
		return
	}
	return os.MkdirAll(r.publishDir, 0755)
}

func (r *Represent) UpdateAssets() error {
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

func (r *Represent) CompileAll() error {
	return filepath.Walk(r.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Clean(path) == filepath.Clean(r.publishDir) {
			// Skip the publish dir
			return filepath.SkipDir
		}
		if !isDoc(path) {
			return nil
		}
		// Derive relative path to static file
		relPath := strings.Replace(path, r.BaseDir, "", 1)
		targetPath := filepath.Join(append([]string{r.publishDir}, filepath.SplitList(relPath)...)...)
		targetPath = regexp.MustCompile(`\.[^.]+$`).ReplaceAllString(targetPath, ".html")
		destFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		return renderDoc(destFile, r.presentPkgDir, path)
	})
}
