package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func (b *Backend) ListDomains() (domains []*backend.Domain, err error) {
	domains = b.domains
	return
}

func (b *Backend) GetDomainByName(name string) (*backend.Domain, error) {
	for _, d := range b.domains {
		if d.Name == name {
			return d, nil
		}
	}
	return nil, errors.New("No such domain")
}
