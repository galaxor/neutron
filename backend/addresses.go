package backend

type AddressesBackend interface {
	// Get a user's address.
	GetAddress(user, id string) (*Address, error)
	// List all addresses owned by a user.
	ListUserAddresses(user string) ([]*Address, error)
	// List all addresses belonging to a domain.
	ListDomainAddresses(domain string) ([]*Address, error)
	// Create a new address.
	InsertAddress(user string, address *Address) (*Address, error)
	// Update an existing address.
	//UpdateAddress(user string, update *AddressUpdate) (*Address, error)
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
