package jsonDB_test

import (
	"bytes"
	"fmt"
	"jsonDB"
	"os"
	"testing"
)

func TestBulkWriter(t *testing.T) {
	fname := "json/j1.json"
	data, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println("can't open " + fname)
		return
	}
	r := bytes.NewReader(data)

	jsonDB.Insert("recipe", r)
}

func BenchmarkBulkWriter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fname := "json/j1.json"
		data, err := os.ReadFile(fname)
		if err != nil {
			return
		}
		r := bytes.NewReader(data)

		jsonDB.Insert("recipe", r)

	}
}
