package models

type URLMapping struct {
	ShortURL string
	LongURL  string
	//creationTime
	//experation date
}

type URLRepository interface {
	Gave(mapping *URLMapping) error
	Get(shortURL string) (*URLMapping, error)
	Delete(shortURL string) error
	Update(mapping *URLMapping)
}
