// Copyright 2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "sevki.org/acme/cmd/header"

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// go get sevki.org/acme/cmd/header

func main() {
	tmpl := template.Must(template.New("header").Parse(hdr))
	dir, _ := os.Getwd()
	rel, _ := filepath.Rel(filepath.Join(os.Getenv("GOPATH"), "src"), dir)
	vars := struct {
		Year    string
		Path    string
		Package string
	}{
		fmt.Sprintf("%d", time.Now().Year()),
		rel,
		filepath.Base(rel),
	}
	if err := tmpl.Execute(os.Stdout, vars); err != nil {
		os.Exit(0)
	}
}

var hdr = `// Copyright {{ .Year }} Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package {{ .Package }} // import "{{ .Path }}"

`

