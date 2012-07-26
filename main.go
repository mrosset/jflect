// Copyright 2012 The jflect Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"unicode"
)

// TODO: write proper Usage and README
var (
	client  = new(http.Client)
	fstruct = flag.String("s", "Foo", "struct name for json object")
	furl    = flag.String("u", "", "url for json input")
	debug   = false
)

func main() {
	flag.Parse()
	if *furl == "" {
		flag.Usage()
		os.Exit(1)
	}
	var v interface{}
	res, err := client.Get(*furl)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		errf := fmt.Errorf("%s %v %s", *furl, res.StatusCode,
			http.StatusText(res.StatusCode))
		log.Fatal(errf)
	}
	err = json.NewDecoder(res.Body).Decode(&v)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	// Open struct
	fmt.Fprintf(buf, "\ntype %s struct {\n", *fstruct)
	b, err := reflect(v)
	if err != nil {
		log.Fatal(err)
	}
	// Write fields to buffer
	buf.Write(b)
	// Close struct
	fmt.Fprintln(buf, "}")
	if debug {
		os.Stdout.WriteString("*********DEBUG***********")
		os.Stdout.Write(buf.Bytes())
		os.Stdout.WriteString("*********DEBUG***********")
	}
	// Pass through gofmt for uniform formatting, and weak syntax check.
	cmd := exec.Command("gofmt")
	cmd.Stdin = buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Field data type
type Field struct {
	name  string
	gtype string
	tag   string
}

// Simplifies Field construction
func NewField(name, gtype string) Field {
	return Field{goField(name), gtype, goTag(name)}
}

// Provides Sorter interface so we can keep field order
type FieldSort []Field

func (s FieldSort) Len() int { return len(s) }

func (s FieldSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s FieldSort) Less(i, j int) bool {
	return s[i].name < s[j].name
}

func reflect(v interface{}) ([]byte, error) {
	var (
		buf = new(bytes.Buffer)
	)
	fields := []Field{}
	switch root := v.(type) {
	case map[string]interface{}:
		for key, val := range root {
			switch j := val.(type) {
			case nil:
				// FIXME: sometimes json service will return nil even though the type is string.
				// go can not convert string to nil and vs versa. Can we assume its a string?
				continue
			case float64:
				fields = append(fields, NewField(key, "int"))
			case map[string]interface{}:
				// If type is map[string]interface{} then we have nested object, Recurse
				fmt.Fprintf(buf, "%s struct {\n", goField(key))
				o, err := reflect(j)
				if err != nil {
					return nil, err
				}
				_, err = buf.Write(o)
				if err != nil {
					return nil, err
				}
				fmt.Fprintln(buf, "}")
			default:
				fields = append(fields, NewField(key, fmt.Sprintf("%T", val)))
			}
		}
	default:
		return nil, fmt.Errorf("%T: unexpected type", root)
	}
	// Sort and write field buffer last to keep order and formatting.
	sort.Sort(FieldSort(fields))
	for _, f := range fields {
		fmt.Fprintf(buf, "%s %s %s\n", f.name, f.gtype, f.tag)
	}
	return buf.Bytes(), nil
}

// Return lower_case json fields to camel case fields.
func goField(jf string) string {
	mkUpper := true
	gf := ""
	for _, c := range jf {
		if mkUpper {
			c = unicode.ToUpper(c)
			mkUpper = false
		}
		if c == '_' {
			mkUpper = true
			continue
		}
		gf += string(c)
	}
	return fmt.Sprintf("%s", gf)
}

// Returns the json tag from a json field.
func goTag(jf string) string {
	return fmt.Sprintf("`json:\"%s\"`", jf)
}
