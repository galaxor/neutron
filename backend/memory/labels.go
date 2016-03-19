package memory

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

type LabelsBackend struct {
	labels map[string][]*backend.Label
}

func (b *LabelsBackend) getLabelIndex(user, id string) (int, error) {
	for i, lbl := range b.labels[user] {
		if lbl.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("No such label")
}

func (b *LabelsBackend) ListLabels(user string) (labels []*backend.Label, err error) {
	labels = b.labels[user]
	return
}

func (b *LabelsBackend) InsertLabel(user string, label *backend.Label) (*backend.Label, error) {
	label.ID = generateId()
	label.Order = len(b.labels[user])
	b.labels[user] = append(b.labels[user], label)
	return label, nil
}

func (b *LabelsBackend) UpdateLabel(user string, update *backend.LabelUpdate) (*backend.Label, error) {
	updated := update.Label

	i, err := b.getLabelIndex(user, updated.ID)
	if err != nil {
		return nil, err
	}

	label := b.labels[user][i]

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

func (b *LabelsBackend) DeleteLabel(user, id string) error {
	i, err := b.getLabelIndex(user, id)
	if err != nil {
		return err
	}

	labels := b.labels[user]
	b.labels[user] = append(labels[:i], labels[i+1:]...)

	return nil
}

func NewLabelsBackend() backend.LabelsBackend {
	return &LabelsBackend{
		labels: map[string][]*backend.Label{},
	}
}
