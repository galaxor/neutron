package backend

type AddressesBackend interface {
	// Get a user's address.
	GetAddress(user, id string) (*Address, error)
	// List all addresses owned by a user.
	ListAddresses(user string) ([]*Address, error)
	// Create a new address.
	InsertAddress(user string, address *Address) (*Address, error)
	// Update an existing address.
	UpdateAddress(user string, update *AddressUpdate) (*Address, error)
	// Delete an address.
	DeleteAddress(user, id string) error
}

// A user's address.
type Address struct {
	ID string
	DomainID string
	Email string
	Send int
	Receive int
	Status int
	Type int
	DisplayName string
	Signature string

	HasKeys int
	Keys []*Keypair
}

// Get this address' email.
func (a *Address) GetEmail() *Email {
	return &Email{
		Address: a.Email,
		Name: a.DisplayName,
	}
}

type AddressUpdate struct {
	Address *Address
	Status bool
	Type bool
	DisplayName bool
	Signature bool
}

func (update *AddressUpdate) Apply(address *Address) {
	updated := update.Address

	if updated.ID != address.ID {
		panic("Cannot apply update on an address with a different ID")
	}

	if update.Status {
		address.Status = updated.Status
	}
	if update.Type {
		address.Type = updated.Type
	}
	if update.DisplayName {
		address.DisplayName = updated.DisplayName
	}
	if update.Signature {
		address.Signature = updated.Signature
	}
}
