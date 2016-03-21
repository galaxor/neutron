package backend

// Stores domains data.
type DomainsBackend interface {
	// List all domains.
	ListDomains() ([]*Domain, error)
	// Get the domain which has the specified name.
	GetDomainByName(name string) (*Domain, error)
}

// A domain name.
type Domain struct {
	ID string
	Name string
}
