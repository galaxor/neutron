package backend

// A contact is an entry in the user's address book.
type Contact struct {
	ID string
	Name string
	Email string
}

// A request to update a contact.
// Fields set to true will be updated with values in Contact.
type ContactUpdate struct {
	Contact *Contact
	Name bool
	Email bool
}
