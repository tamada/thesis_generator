package main

import "testing"

func TestDo(t *testing.T) {
	goMain([]string{"thesis_generator", "../../testdata/sample.json", "-o", "../../dest.zip"})
}
