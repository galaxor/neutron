package util

import (
	"github.com/emersion/neutron/backend"
)

type EventedLabels struct {
	backend.LabelsBackend
	events backend.EventsBackend
}

func (b *EventedLabels) InsertLabel(user string, label *backend.Label) (*backend.Label, error) {
	label, err := b.LabelsBackend.InsertLabel(user, label)

	if err == nil {
		event := backend.NewLabelDeltaEvent(label.ID, backend.EventCreate, label)
		b.events.InsertEvent(user, event)
	}

	return label, err
}

func (b *EventedLabels) UpdateLabel(user string, update *backend.LabelUpdate) (*backend.Label, error) {
	label, err := b.LabelsBackend.UpdateLabel(user, update)

	if err == nil {
		event := backend.NewLabelDeltaEvent(label.ID, backend.EventUpdate, label)
		b.events.InsertEvent(user, event)
	}

	return label, err
}

func (b *EventedLabels) DeleteLabel(user, id string) error {
	err := b.LabelsBackend.DeleteLabel(user, id)

	if err == nil {
		event := backend.NewLabelDeltaEvent(id, backend.EventDelete, nil)
		b.events.InsertEvent(user, event)
	}

	return err
}

func NewEventedLabels(bkd backend.LabelsBackend, events backend.EventsBackend) backend.LabelsBackend {
	return &EventedLabels{
		LabelsBackend: bkd,
		events: events,
	}
}
