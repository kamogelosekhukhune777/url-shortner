package main

import (
	"testing"

	"github.com/kamogelosekhukhune777/url-shortner/data"
	"github.com/stretchr/testify/assert"
)

func TestSimpleRedirection(t *testing.T) {
	store := data.NewInmemLinkStore()

	want := `http://go.dev/`
	got, err := store.GetLinkByShort(want)
	assert.NoError(t, err)
	assert.Equal(t, want, got.FullURL)
}
