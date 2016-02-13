package models

type Phone struct {
	Type      string `gorethink:"type" json:"type"`
	Number    string `gorethink:"number" json:"number"`
	Extension string `gorethink:"extension" json:"extension,omitempty"`
}

type Contact struct {
	ID        string   `gorethink:"id,omitempty" json:"-"`
	Name      string   `gorethink:"name" json:"name"`
	FirstName string   `gorethink:"firstName" json:"firstName"`
	LastName  string   `gorethink:"lastName" json:"lastName"`
	Synonyms  []string `gorethink:"synonyms" json:"synonyms"`
	Position  string   `gorethink:"position" json:"position"`
	Phones    []Phone  `gorethink:"phones" json:"phones"`
	Emails    []string `gorethink:"emails" json:"emails"`
	Notes     string   `gorethink:"notes" json:"notes"`
	Related   []string `gorethink:"related" json:"related"`
}

func (c Contact) GetName() string {
	return "contacts"
}

func (c *Contact) SetID(ID string) error {
	c.ID = ID
	return nil
}

func (c Contact) GetID() string {
	return c.ID
}
