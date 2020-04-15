package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"9fans.net/go/acme"
	"github.com/google/subcommands"
	"github.com/hexdigest/gounit"
)

type genType int

const (
	tests genType = iota
	mocks
)

type generator struct{ genType }

func (g *generator) Name() string {
	return g.genType.String()
}
func (g *generator) Synopsis() string {
	return fmt.Sprintf("generate %s from a file from an acme window", g.genType)
}
func (*generator) Usage() string {
	return `tests`
}

func (g *generator) SetFlags(f *flag.FlagSet) {
}
func srcToTest(src string) string {
	return strings.Replace(src, ".go", "_test.go", -1)
}

func (g *generator) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	winid, ok := ctx.Value(winIDKey).(int)
	if !ok || winid == 0 {
		usage()
		return subcommands.ExitUsageError
	}

	file, ok := ctx.Value(fileKey).(string)
	if !ok || !strings.HasSuffix(file, ".go") {
		usage()
		return subcommands.ExitUsageError
	}
	testFile := srcToTest(file)
	w, err := acme.Open(winid, nil)
	if err != nil {
		panic(err)
	}

	var testBody io.ReadWriter
	testReader := testBody
	testWindow := findWindow(testFile)
	if testWindow != nil {
		testBody = bodyReadWriter{testWindow}
		testReader = testBody
	} else {
		testWindow, err = acme.New()
		if err != nil {
			panic(err)
		}
		testWindow.Name(testFile)
		testBody = bodyReadWriter{testWindow}
		testReader = nil
	}
	opts := gounit.Options{
		UseStdin:   true,
		UseJSON:    false,
		All:        true,
		InputFile:  file,
		OutputFile: srcToTest(file),
		UseStdout:  true,
	}
	switch g.genType {
	case tests:
		opts.Template = testify
		opts.TemplateName = "testify"
	case mocks:
		opts.Template = minimock
		opts.Template = "minimock"
	}
GEN:
	gen, err := gounit.NewGenerator(opts, bodyReadWriter{w}, testReader)
	if err == gounit.ErrFuncNotFound {
		testReader = nil
		goto GEN
	} else if err != nil {
		panic(err)
	}
	_, err = testWindow.Seek("body", 0, 0)
	gen.Write(testBody)
	return subcommands.ExitSuccess
}

const minimock = `{{$func := .Func}}

func {{ $func.TestName }}(t *testing.T) {
	{{- if (gt $func.NumParams 0) }}
		type args struct {
			{{ range $param := params $func }}
				{{- $param}}
			{{ end }}
		}
	{{ end -}}
	tests := []struct {
		name string
		{{- if $func.IsMethod }}
			init func(t minimock.Tester) {{ ast $func.ReceiverType }}
			inspect func(r {{ ast $func.ReceiverType }}, t *testing.T) //inspects {{ ast $func.ReceiverType }} after execution of {{$func.Name}}
		{{ end }}
		{{- if (gt $func.NumParams 0) }}
			args func(t minimock.Tester) args
		{{ end }}
		{{ range $result := results $func}}
			{{ want $result -}}
		{{ end }}
		{{- if $func.ReturnsError }}
			wantErr bool
			inspectErr func (err error, t *testing.T) //use for more precise error evaluation
		{{ end -}}
	}{
		{{- if eq .Comment "" }}
			//TODO: Add test cases
		{{else}}
			//{{ .Comment }}
		{{end -}}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		  mc := minimock.NewController(t)
		  defer mc.Wait(time.Second)

			{{if (gt $func.NumParams 0) }}
				tArgs := tt.args(mc)
			{{ end}}
			{{- if $func.IsMethod -}}
				receiver := tt.init(mc)

				{{ if (gt $func.NumResults 0) }}{{ join $func.ResultsNames ", " }} := {{end}}receiver.{{$func.Name}}(
					{{- range $i, $pn := $func.ParamsNames }}
						{{- if not (eq $i 0)}},{{end}}tArgs.{{ $pn }}{{ end }})

				if tt.inspect != nil {
					tt.inspect(receiver, t)
				}
			{{ else }}
				{{ if (gt $func.NumResults 0) }}{{ join $func.ResultsNames ", " }} := {{end}}{{$func.Name}}(
					{{- range $i, $pn := $func.ParamsNames }}
						{{- if not (eq $i 0)}},{{end}}tArgs.{{ $pn }}{{ end }})
			{{end}}
			{{ range $result := $func.ResultsNames }}
				{{ if (eq $result "err") }}
          if tt.wantErr {
            if assert.Error(t, err) && tt.inspectErr!= nil {
						  tt.inspectErr(err, t)
            }
					} else {
					  assert.NoError(t, err)
					}
					
				{{ else }}
				  assert.Equal(t, tt.{{ want $result }}, {{ $result }}, "{{ receiver $func }}{{ $func.Name }} returned unexpected result")
				{{end -}}
			{{end -}}
		})
	}
}`
const testify = `{{$func := .Func}}

func {{ $func.TestName }}(t *testing.T) {
	{{- if (gt $func.NumParams 0) }}
		type args struct {
			{{ range $param := params $func }}
				{{- $param}}
			{{ end }}
		}
	{{ end -}}
	tests := []struct {
		name string
		{{- if $func.IsMethod }}
			init func(t *testing.T) {{ ast $func.ReceiverType }}
			inspect func(r {{ ast $func.ReceiverType }}, t *testing.T) 
		{{ end }}
		{{- if (gt $func.NumParams 0) }}
			args func(t *testing.T) args
		{{ end }}
		{{ range $result := results $func}}
			{{ want $result -}}
		{{ end }}
		{{- if $func.ReturnsError }}
			wantErr bool
			inspectErr func (err error, t *testing.T)
		{{ end -}}
	}{
		{{- if eq .Comment "" }}
			//TODO: Add test cases
		{{else}}
			//{{ .Comment }}
		{{end -}}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{{- if (gt $func.NumParams 0) }}
				tArgs := tt.args(t)
			{{ end -}}
			{{ if $func.IsMethod }}
				receiver := tt.init(t)
				{{ if (gt $func.NumResults 0) }}{{ join $func.ResultsNames ", " }} := {{end}}receiver.{{$func.Name}}(
					{{- range $i, $pn := $func.ParamsNames }}
						{{- if not (eq $i 0)}},{{end}}tArgs.{{ $pn }}{{ end }})

				if tt.inspect != nil {
					tt.inspect(receiver, t)
				}
			{{ else }}
				{{ if (gt $func.NumResults 0) }}{{ join $func.ResultsNames ", " }} := {{end}}{{$func.Name}}(
					{{- range $i, $pn := $func.ParamsNames }}
						{{- if not (eq $i 0)}},{{end}}tArgs.{{ $pn }}{{ end }})
			{{end}}
			{{ range $result := $func.ResultsNames }}
				{{ if (eq $result "err") }}
				if tt.wantErr {
				  require.Error(t, err)
          if tt.inspectErr!= nil {
						tt.inspectErr(err, t)
					}
				}
				{{ else }}
				  assert.Equal(t, tt.{{ want $result }}, {{ $result }})
				{{end -}}
			{{end -}}
		})
	}
}
`
