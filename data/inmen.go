package data

import "fmt"

type InmemLinkStore struct {
	data []Link
}

func NewInmemLinkStore() *InmemLinkStore {
	l1 := NewLink("https://go.dev")
	l1.Short = "go"

	l2 := NewLink("https://example.com")
	l2.Short = "example"

	l3 := NewLink("https://yahoo.com")
	l3.Short = "yahoo"

	return &InmemLinkStore{
		data: []Link{l1, l2, l3},
	}
}

func (store *InmemLinkStore) Save(link Link) error {
	for i, l := range store.data {
		if l.ID == link.ID {
			store.data[i] = link
			return nil
		}
	}
	store.data = append(store.data, link)
	return nil
}

func (store *InmemLinkStore) GetLinkByID(id string) (*Link, error) {
	for _, l := range store.data {
		if l.ID == id {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("link '%s' not found", id)
}

func (store *InmemLinkStore) GetLinkByShort(short string) (*Link, error) {
	for _, l := range store.data {
		if l.Short == short {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("link '%s' not found by short", short)
}
