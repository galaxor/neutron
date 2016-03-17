package backend

type Contact struct {
	ID string
	Name string
	Email string
}

type ContactUpdate struct {
	Contact *Contact
	Name bool
	Email bool
}
