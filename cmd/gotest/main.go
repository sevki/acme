// Copyright 2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "sevki.org/acme/cmd/gotest"

import (
	"os"
	"text/template"
)

func main() {
	tmpl := template.Must(template.New("header").Parse(hdr))

	vars := struct {
		Name string
	}{
		os.Args[1],
	}
	if err := tmpl.Execute(os.Stdout, vars); err != nil {
		os.Exit(0)
	}
}

var hdr = `func Test{{.Name}}(t *testing.T) {
	
}`
