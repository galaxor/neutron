package backend

// Stores domains data.
type DomainsBackend interface {
	// List all domains.
	ListDomains() ([]*Domain, error)
	// Get a domain.
	GetDomain(id string) (*Domain, error)
	// Get the domain which has the specified name.
	GetDomainByName(name string) (*Domain, error)
	// Insert a new domain.
	InsertDomain(domain *Domain) (*Domain, error)
}

// A domain name.
type Domain struct {
	ID string
	DomainName string

	State int
	VerifyState int
	MxState int
	SpfState int
	DkimState int
	DmarcState int

	Addresses []*Address
}
