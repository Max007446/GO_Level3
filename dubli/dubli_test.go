package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testwalk_single(t *testing.T, path *string) {
	assert.DirExists(t, *path, "norm")
}
func testwalk_multi(t *testing.T, path *string) {
	assert.FileExists(t, *path, "norm")
}
