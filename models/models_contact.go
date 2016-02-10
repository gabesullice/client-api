package models

type Phone struct {
	Type      string `json:"type"`
	Number    string `json:"number"`
	Extension string `json:"extension,omitempty"`
}

type Contact struct {
	ID        string   `json:"-"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Position  string   `json:"position"`
	Phones    []Phone  `json:"phones"`
	Emails    []string `json:"emails"`
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
