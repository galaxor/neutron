package backend

// A user.
type User struct {
	ID string
	Name string
	NotificationEmail string
	Signature string
	NumMessagePerPage int
	UsedSpace int
	Notify int
	AutoSaveContacts int
	Language string
	LogAuth int
	ComposerMode int
	MessageButtons int
	ShowImages int
	ViewMode int
	ViewLayout int
	SwipeLeft int
	SwipeRight int
	Theme string
	Currency string
	Credit int
	DisplayName string
	MaxSpace int
	MaxUpload int
	Role int
	Private int
	Subscribed int
	Deliquent int
	Addresses []*Address
	PublicKey string
	EncPrivateKey string
}

func (u *User) GetMainAddress() *Address {
	for _, addr := range u.Addresses {
		if addr.Send == 1 { // 1 is main address, 2 is secondary address
			return addr
		}
	}
	return nil
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

// A request to update a user.
// Fields set to true will be updated with values in User.
type UserUpdate struct {
	User *User
	DisplayName bool
	Signature bool
	AutoSaveContacts bool
	ShowImages bool
}
