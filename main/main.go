// Copyright 2012 The jflect Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"github.com/mathew-bowersox/jflect"
	glog "log"
	"os"
)

// TODO: write proper Usage and README
var (
	log               = glog.New(os.Stderr, "", glog.Lshortfile)
	fstruct           = flag.String("s", "Foo", "struct name for json object")
	debug             = false
	ErrNotValidSyntax = errors.New("Json reflection is not valid Go syntax")
)

func main() {
	flag.Parse()
	err := generate.Generate(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}


