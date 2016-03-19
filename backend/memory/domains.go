package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

type DomainsBackend struct {
	domains []*backend.Domain
}

func (b *DomainsBackend) ListDomains() (domains []*backend.Domain, err error) {
	domains = b.domains
	return
}

func (b *DomainsBackend) GetDomainByName(name string) (*backend.Domain, error) {
	for _, d := range b.domains {
		if d.Name == name {
			return d, nil
		}
	}
	return nil, errors.New("No such domain")
}

func NewDomainsBackend() backend.DomainsBackend {
	return &DomainsBackend{}
}
