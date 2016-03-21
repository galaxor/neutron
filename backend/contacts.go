package backend

// Stores contacts data.
type ContactsBackend interface {
	// List all user's contacts.
	ListContacts(user string) ([]*Contact, error)
	// Insert a new contact.
	InsertContact(user string, contact *Contact) (*Contact, error)
	// Update an existing contact.
	UpdateContact(user string, update *ContactUpdate) (*Contact, error)
	// Delete a contact.
	DeleteContact(user, id string) error
	// Delete all contacts of a specific user.
	DeleteAllContacts(user string) error
}

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

// Apply this update on a contact.
func (update *ContactUpdate) Apply(contact *Contact) {
	updated := update.Contact

	if updated.ID != contact.ID {
		panic("Cannot apply update on a contact with a different ID")
	}

	if update.Name {
		contact.Name = updated.Name
	}
	if update.Email {
		contact.Email = updated.Email
	}
}
