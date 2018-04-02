// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Acmefmt watches acme for a variety files being written.
// Each time a file is written, acmefmt finds the relevant formatter
// and formats and writes back the file.
//
// Registered formatters are:
/*
	// Dockerfile -> dockfmt
	// BUILD -> buildifier

	// goimports
	".go": goimports{},

	// clang-format
	".c":   clangfmt{},
	".h":   clangfmt{},
	".cc":  clangfmt{},
	".cpp": clangfmt{},
	".hpp": clangfmt{},

	// prettier
	".js":   prettier{},
	".json": prettier{},
	".md":   prettier{},

	// rustfmt
	".rs": rustfmt{},

	// buildifier
	".bzl": buildifier{},
*/
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

var (
	fmtrs = map[string]formatter{
		// goimports
		".go": goimports{},

		// clang-format
		".c":   clangfmt{},
		".h":   clangfmt{},
		".cc":  clangfmt{},
		".cpp": clangfmt{},
		".hpp": clangfmt{},

		// prettier
		".js":   prettier{},
		".json": prettier{},
		".md":   prettier{},

		// rustfmt
		".rs": rustfmt{},

		// buildifier
		".bzl": buildifier{},
	}
)

func main() {
	flag.Parse()

	l, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}
		filter(event)
	}
}

type formatter interface {
	format(string) ([]byte, error)
}

func filter(e acme.LogEvent) {
	if e.Op != "put" {
		return
	}

	_, file := path.Split(e.Name)
	ext := path.Ext(file)

	switch file {
	case "Dockerfile":
		reformat(e.ID, e.Name, dockfmt{})
	case "BUILD", "BUCK":
		reformat(e.ID, e.Name, buildifier{})
	default:
		if fmtr, ok := fmtrs[ext]; ok {
			reformat(e.ID, e.Name, fmtr)
		}
	}
}

func reformat(id int, name string, formatter formatter) {
	w, err := acme.Open(id, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()

	old, err := ioutil.ReadFile(name)
	if err != nil {
		//log.Print(err)
		return
	}
	new, err := formatter.format(name)
	if bytes.Equal(old, new) {
		return
	}

	f, err := ioutil.TempFile("", "acmefmt")
	if err != nil {
		log.Print(err)
		return
	}
	if _, err := f.Write(new); err != nil {
		log.Print(err)
		return
	}
	tmp := f.Name()
	f.Close()
	defer os.Remove(tmp)

	diff, _ := exec.Command("9", "diff", name, tmp).CombinedOutput()

	w.Write("ctl", []byte("mark"))
	w.Write("ctl", []byte("nomark"))
	diffLines := strings.Split(string(diff), "\n")
	for i := len(diffLines) - 1; i >= 0; i-- {
		line := diffLines[i]
		if line == "" {
			continue
		}
		if line[0] == '<' || line[0] == '-' || line[0] == '>' {
			continue
		}
		j := 0
		for j < len(line) && line[j] != 'a' && line[j] != 'c' && line[j] != 'd' {
			j++
		}
		if j >= len(line) {
			log.Printf("cannot parse diff line: %q", line)
			break
		}
		oldStart, oldEnd := parseSpan(line[:j])
		newStart, newEnd := parseSpan(line[j+1:])
		if oldStart == 0 || newStart == 0 {
			continue
		}
		switch line[j] {
		case 'a':
			err := w.Addr("%d+#0", oldStart)
			if err != nil {
				log.Print(err)
				break
			}
			w.Write("data", findLines(new, newStart, newEnd))
		case 'c':
			err := w.Addr("%d,%d", oldStart, oldEnd)
			if err != nil {
				log.Print(err)
				break
			}
			w.Write("data", findLines(new, newStart, newEnd))
		case 'd':
			err := w.Addr("%d,%d", oldStart, oldEnd)
			if err != nil {
				log.Print(err)
				break
			}
			w.Write("data", nil)
		}
	}
}

func parseSpan(text string) (start, end int) {
	i := strings.Index(text, ",")
	if i < 0 {
		n, err := strconv.Atoi(text)
		if err != nil {
			log.Printf("cannot parse span %q", text)
			return 0, 0
		}
		return n, n
	}
	start, err1 := strconv.Atoi(text[:i])
	end, err2 := strconv.Atoi(text[i+1:])
	if err1 != nil || err2 != nil {
		log.Printf("cannot parse span %q", text)
		return 0, 0
	}
	return start, end
}

func findLines(text []byte, start, end int) []byte {
	i := 0

	start--
	for ; i < len(text) && start > 0; i++ {
		if text[i] == '\n' {
			start--
			end--
		}
	}
	startByte := i
	for ; i < len(text) && end > 0; i++ {
		if text[i] == '\n' {
			end--
		}
	}
	endByte := i
	return text[startByte:endByte]
}
