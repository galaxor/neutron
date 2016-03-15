package memory

import (
	"github.com/emersion/neutron/backend"
)

func (b *Backend) ListLabels(user string) (labels []*backend.Label, err error) {
	labels = b.data[user].labels
	return
}
