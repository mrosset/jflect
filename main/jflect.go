// Copyright 2012 The jflect Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mathew-bowersox/jflect"
	glog "log"
	"os"
)

// TODO: write proper Usage and README

func main() {
	err := generate.Generate(os.Stdin, os.Stdout)
	if err != nil {
		glog.Fatal(err)
	}
}


