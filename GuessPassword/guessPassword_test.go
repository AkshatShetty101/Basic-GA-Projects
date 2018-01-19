package main

import (
	"strings"
	"testing"
)

func TestGenerateParent(t *testing.T) {
	geneset := " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!."
	sliceGeneset := strings.Split(geneset, "")
	// target := "Hello World!"
	str := generateParent(sliceGeneset, 12)
	if len(str) != 12 {
		t.Errorf("Expected length 12 but got %v", len(str))
	}
}
