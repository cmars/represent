// Copyright 2013 Casey Marshall <casey.marshall@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	. "github.com/cmars/represent/pkg/represent"
)

var src *string = flag.String("src", "", "Source path containing Present files and referenced content")
var publish *string = flag.String("publish", "", "Publish path to create static HTML pages and assets")
var base *string = flag.String("base", "", "Base path for slide template and static resources")

func die(err error, v ...interface{}) {
	log.Println(append(v, err)...)
}

func main() {
	flag.Parse()
	represent, err := NewRepresent(*src, *publish, *base)
	if err != nil {
		die(err)
	}
	err = represent.Publish()
	if err != nil {
		die(err)
	}
}
