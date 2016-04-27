// Copyright 2016 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmnt // import "sevki.org/acme/cmd/cmnt"
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		scanned := scanner.Text()
		if trimmed := strings.TrimLeft(scanned, "//"); trimmed == scanned {
			fmt.Fprintf(os.Stdout, "//%s\n", scanned)
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", trimmed)
		}
	}

}
