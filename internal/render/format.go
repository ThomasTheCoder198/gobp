package render

import (
	"fmt"
	"go/format"
)

// gofmt runs go/format.Source over the rendered bytes.
//
// On failure the error includes a snippet of the offending output so a
// template bug is debuggable from the test log.
func gofmt(target string, src []byte) ([]byte, error) {
	out, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("gofmt %s: %w\n--- rendered output ---\n%s\n--- end ---", target, err, src)
	}
	return out, nil
}
