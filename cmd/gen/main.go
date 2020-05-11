package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/markbates/uad"

	"github.com/gobuffalo/flect"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	f, err := os.Open("plugins.json")
	if err != nil {
		return err
	}
	defer f.Close()

	plugs := map[string]uad.Plugin{}
	if err := json.NewDecoder(f).Decode(&plugs); err != nil {
		return err
	}
	fmt.Println(">>>TODO cmd/gen/main.go:13: plugs ", plugs)

	cats := map[string][]uad.Plugin{}

	for _, p := range plugs {
		cats[p.Category] = append(cats[p.Category], p)
	}

	os.RemoveAll("plugins")
	os.MkdirAll("plugins", 0755)
	for c, plugs := range cats {
		bb := &bytes.Buffer{}
		bb.WriteString("package plugins\n\n")
		bb.WriteString(fmt.Sprintf("type %s struct{}\n\n", flect.Pascalize(c)))
		for _, p := range plugs {
			bb.WriteString(p.String())
		}
		f, err := os.Create(filepath.Join("plugins", flect.Underscore(c)+".go"))
		if err != nil {
			return err
		}
		defer f.Close()
		// io.Copy(os.Stdout, bb)
		io.Copy(f, bb)
	}

	return nil
}
