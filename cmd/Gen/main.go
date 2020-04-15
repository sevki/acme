// Copyright 2016-2020 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"9fans.net/go/acme"
	"github.com/google/subcommands"
)

type acmeKey string

const (
	winIDKey acmeKey = "_winid"
	fileKey  acmeKey = "_file"
	fontKey  acmeKey = "_font"
)

var ()

type bodyReadWriter struct{ *acme.Win }

func (r bodyReadWriter) Write(b []byte) (n int, err error) { return r.Win.Write("body", b) }
func (r bodyReadWriter) Read(b []byte) (n int, err error)  { return r.Win.Read("body", b) }

func findWindow(name string) *acme.Win {
	re := regexp.MustCompile(name)
	windows, err := acme.Windows()
	if err != nil {
		panic(err)
	}

	for _, w := range windows {
		if re.MatchString(w.Name) {
			if win, err := acme.Open(w.ID, nil); err == nil {
				return win
			}

		}
	}

	return nil
}

func usage() {
	fmt.Println(`Gen us supposed to be invoed within acme`)
	os.Exit(33)
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&generator{mocks}, "gounit")
	subcommands.Register(&generator{tests}, "gounit")

	winid, _ := strconv.Atoi(os.Getenv("winid"))
	file := os.Getenv("%")
	font := os.Getenv("font")

	flag.Parse()
	ctx := context.Background()
	ctx = context.WithValue(ctx, winIDKey, winid)
	ctx = context.WithValue(ctx, fileKey, file)
	ctx = context.WithValue(ctx, fontKey, font)

	os.Exit(int(subcommands.Execute(ctx)))

}
