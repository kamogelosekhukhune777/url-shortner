package main

import (
	"testing"

	"github.com/kamogelosekhukhune777/url-shortner/data"
)

func TestSimpleRedirection(t *testing.T) {
	store := data.NewInmemLinkStore()

	want := `http://go.dev/`
	got, err := store.GetLinkByShort(want)
	assert.No(t, err)
}
