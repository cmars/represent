// Copyright 2013 Casey Marshall <casey.marshall@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	. "github.com/cmars/represent"
	"os"
)

var baseDir *string = flag.String("basedir", "", "Base directory of static site")

func main() {
	flag.Parse()
	if *baseDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		baseDir = &wd
	}
	represent, err := NewRepresent(*baseDir)
	if err != nil {
		panic(err)
	}
	// Create publish directory
	represent.RequirePublishDir()
	// Copy present static files to publish directory
	err = represent.UpdateAssets()
	if err != nil {
		panic(err)
	}
	// For each present file, render to html in publish directory
	err = represent.CompileAll()
	if err != nil {
		panic(err)
	}
}
