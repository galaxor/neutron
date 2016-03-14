package memory

import (
	"github.com/emersion/neutron/backend"
)

type Backend struct {}

func New() backend.Backend {
	return &Backend{}
}
