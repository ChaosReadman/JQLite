package bulkWriter_test

import (
	"bulkWriter"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestBulkWriter(t *testing.T) {
	fname := "json/j1.json"
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("can't open " + fname)
		return
	}
	r := io.Reader(file)

	bulkWriter.WriteJson("recipe", r)
}
