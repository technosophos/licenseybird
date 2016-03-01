package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const LicenseBlock = `Copyright {{.Year}} {{.Owners}} All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`

const DefaultOwners = "The Kubernetes Authors"

var lic string

func main() {
	t := template.Must(template.New("t").Parse(LicenseBlock))
	y := time.Now().Year()
	o := os.Getenv("LICENSEYBIRD_OWNERS")
	if o == "" {
		o = DefaultOwners
	}
	v := map[string]interface{}{
		"Year":   y,
		"Owners": o,
	}

	var b bytes.Buffer
	if err := t.Execute(&b, v); err != nil {
		fmt.Printf("Template failed: %s", err)
		os.Exit(127)
	}
	lic = b.String()

	exit := 0
	tmp := bytes.NewBuffer(nil)
	for _, fname := range os.Args[1:] {

		fi, err := os.Stat(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping file %s: %s", fname, err)
			continue
		}

		fmt.Printf("Adding license to %s\n", fname)
		if err := addLicense(fname, tmp); err != nil {
			fmt.Printf("Skipped %s: %s", fname, err)
			exit++
			continue
		}

		// If we get here, we can overwrite the old file.
		ioutil.WriteFile(fname, tmp.Bytes(), fi.Mode())
		tmp.Reset()

	}
	os.Exit(exit)
}

func addLicense(fname string, out io.Writer) error {

	in, err := os.Open(fname)
	if err != nil {
		return err
	}

	// Special handling for certain files.
	b := filepath.Base(fname)
	switch b {
	case "Makefile", "Dockerfile":
		if err := hashPre(out); err != nil {
			return err
		}
		_, err := io.Copy(out, in)
		return err
	}

	ext := filepath.Ext(fname)
	switch ext {
	case ".py":
		fmt.Fprintln(out, "#!/usr/bin/env python\n")
		hashPre(out)
	case ".sh", ".bash":
		hashPre(out)
	case ".md":
		fmt.Fprintln(os.Stderr, "Markdown files do not need license blocks.")
		return nil
	case ".go":
		fmt.Println("golang")
		slashStarPre(out)
	default:
		return fmt.Errorf("Format %s not supported", ext)
	}

	_, err = io.Copy(out, in)
	return err
}

func hashPre(out io.Writer) error {
	return linePrefix(lic, "#", out)
}

func slashPre(out io.Writer) error {
	return linePrefix(lic, "//", out)
}

func slashStarPre(out io.Writer) {
	fmt.Fprintln(out, "/*")
	fmt.Fprint(out, lic)
	fmt.Fprintln(out, "*/")
}

// linePrefix takes a multi-line source text and prepends the prefix to each line.
func linePrefix(source, prefix string, out io.Writer) error {
	p := []byte(prefix)
	b := bytes.NewBuffer([]byte(source))
	eol := []byte("\n")
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		out.Write(p)
		d := scanner.Bytes()
		if len(d) > 1 {
			out.Write([]byte(" "))
		}
		out.Write(d)
		out.Write(eol)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
