package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
)

var (
	url = "https://api.github.com/repos/str1ngs/gotimer"
)

func TestReflect(t *testing.T) {
	os.Args = append(os.Args, "-u", url)
	main()
}

func decode(v *interface{}, url string) (err error) {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s %v %s", url, res.StatusCode,
			http.StatusText(res.StatusCode))
	}
	return json.NewDecoder(res.Body).Decode(&v)
}
