package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type goimports struct{}

func (goimports) format(name string) ([]byte, error) {

	new, err := exec.Command("goimports", name).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "goimports %s: %v\n%s", name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return nil, err
	}
	return new, nil
}

type prettier struct{}

func (prettier) format(name string) ([]byte, error) {

	new, err := exec.Command("prettier", name).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "prettier %s: %v\n%s", name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return nil, err
	}
	return new, nil
}

type clangfmt struct{}

func (clangfmt) format(name string) ([]byte, error) {

	new, err := exec.Command("clang-format", name).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "clang-format %s: %v\n%s", name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return nil, err
	}
	return new, nil
}

type rustfmt struct{}

func (rustfmt) format(name string) ([]byte, error) {

	new, err := exec.Command("rustfmt", "--color", "never", "--write-mode", "plain", name).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "rustfmt %s: %v\n%s", name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return nil, err
	}
	return new, nil
}

type dockfmt struct{}

func (dockfmt) format(name string) ([]byte, error) {

	new, err := exec.Command("dockfmt", "fmt", name).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "dockfmt %s: %v\n%s", name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return nil, err
	}
	return new, nil
}

type buildifier struct{}

func (buildifier) format(name string) ([]byte, error) {

	new, err := exec.Command("buildifier", "-mode=print_if_changed", name).CombinedOutput()
	if err != nil {
		if strings.Contains(string(new), "fatal error") {
			fmt.Fprintf(os.Stderr, "dockfmt %s: %v\n%s", name, err, new)
		} else {
			fmt.Fprintf(os.Stderr, "%s", new)
		}
		return nil, err
	}
	if len(new) <= 0 {
		return ioutil.ReadFile(name)
	}
	return new, nil
}