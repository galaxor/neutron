package memory

import (
	"github.com/emersion/neutron/backend"
)

func (b *Backend) GetLabels(user string) (labels []*backend.Label, err error) {
	labels = []*backend.Label{
		&backend.Label{
			ID: "label_id",
			Name: "Hey!",
			Color: "#7272a7",
			Display: 1,
			Order: 1,
		},
	}
	return
}
