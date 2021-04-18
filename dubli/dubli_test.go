package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testwalk_single(t *testing.T, path *string) { // is unused (deadcode)
	assert.DirExists(t, *path, "norm")
}
func testwalk_multi(t *testing.T, path *string) { // is unused (deadcode)
	assert.FileExists(t, *path, "norm")
}
