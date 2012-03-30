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
	"unicode"
)

var (
	client  = new(http.Client)
	fstruct = flag.String("s", "User", "struct name for json object")
	furl    = flag.String("u", "https://api.github.com/users/str1ngs", "url for json input")
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
	if err = reflect(os.Stdout, v, *fstruct); err != nil {
		log.Fatal(err)
	}
}

func reflect(w io.Writer, i interface{}, strct string) (err error) {
	bb := new(bytes.Buffer)
	switch i := i.(type) {
	case map[string]interface{}:
		fmt.Fprintf(bb, "type %s struct {\n", strct)
		for key, val := range i {
			if len(key) == 0 {
				return fmt.Errorf("len or map key is 0")
			}
			gotype := fmt.Sprintf("%T", val)
			switch gotype {
			case "<nil>":
				gotype = "nil"
			case "float64":
				gotype = "int"
			}
			mkUpper := true
			field := ""
			for _, c := range key {
				if mkUpper {
					c = unicode.ToUpper(c)
					mkUpper = false
				}
				if c == '_' {
					mkUpper = true
					continue
				}
				field += string(c)
			}
			fmt.Fprintf(bb, "%s %s `json:\"%s\"`\n", field, gotype, key)
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
