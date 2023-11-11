package data

var newID, _ = nonoid.Standard(25)
var newShort, _ = nonoid.CustomASCII("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuv", 8)

// all of our models
type Link struct {
	ID      string `json:"id"`
	Short   string `json:"short"`
	FullURL string `json:"full_URL"`
	Active  bool   `json:"active"`
}

func NewLink(full string) Link {
	return Link{
		ID:      newID,
		Short:   newShort,
		FullURL: full,
		Active:  true,
	}
}

type LinkStorer interface {
	Save(link Link) error
	GetLinkByID(id string) (*Link, error)
	GetLinkByShort(short string) (*Link, error)
}
