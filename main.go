package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	client   = new(http.Client)
	fstruct  = flag.String("s", "Foo", "struct name for json object")
	fpackage = flag.String("p", "main", "package name")
	furl     = flag.String("u", "", "url for json input")
)

func main() {
	flag.Parse()
	if *furl == "" {
		flag.Usage()
		os.Exit(1)
	}
	//var tjson = []byte(`{ "name": "Joe", "age": 25 }`)
	var v interface{}
	res, err := client.Get(*furl)
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(res.Body).Decode(&v)
	if err != nil {
		log.Fatal(err)
	}
	if err = Reflect(os.Stdout, v, *fpackage, *fstruct); err != nil {
		log.Fatal(err)
	}
}

func Reflect(w io.Writer, i interface{}, pkg, strct string) (err error) {
	bb := new(bytes.Buffer)
	switch i := i.(type) {
	case map[string]interface{}:
		fmt.Fprintf(bb, "package %s\n\n", pkg)
		fmt.Fprintf(bb, "type %s struct {\n", strct)
		for key, val := range i {
			vstr := fmt.Sprintf("%T", val)
			if vstr == "<nil>" {
				vstr = "nil"
			}
			fmt.Fprintf(bb, "%s %s\n", key, vstr)
		}
		fmt.Fprintln(bb, "}")
		cmd := exec.Command("gofmt")
		cmd.Stdin = bb
		cmd.Stdout = w
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unexpected type")
	}
	return nil
}
