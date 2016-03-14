package backend

type Contact struct {
	ID string
	Name string
	Email string
}

func GetContacts(id string) (contacts []*Contact, err error) {
	contacts = []*Contact{
		&Contact{
			ID: "contact_id",
			Name: "Myself :)",
			Email: "neutron@example.org",
		},
	}
	return
}
