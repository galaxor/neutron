package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func (b *Backend) getLabelIndex(user, id string) (int, error) {
	for i, lbl := range b.data[user].labels {
		if lbl.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("No such label")
}

func (b *Backend) ListLabels(user string) (labels []*backend.Label, err error) {
	labels = b.data[user].labels
	return
}

func (b *Backend) InsertLabel(user string, label *backend.Label) (*backend.Label, error) {
	label.ID = generateId()
	label.Order = len(b.data[user].labels)
	b.data[user].labels = append(b.data[user].labels, label)
	return label, nil
}

func (b *Backend) UpdateLabel(user string, update *backend.LabelUpdate) (*backend.Label, error) {
	updated := update.Label

	i, err := b.getLabelIndex(user, updated.ID)
	if err != nil {
		return nil, err
	}

	label := b.data[user].labels[i]

	if update.Name {
		label.Name = updated.Name
	}
	if update.Color {
		label.Color = updated.Color
	}
	if update.Display {
		label.Display = updated.Display
	}
	if update.Order {
		label.Order = updated.Order
	}

	return label, nil
}

func (b *Backend) DeleteLabel(user, id string) error {
	i, err := b.getLabelIndex(user, id)
	if err != nil {
		return err
	}

	labels := b.data[user].labels
	b.data[user].labels = append(labels[:i], labels[i+1:]...)

	return nil
}
