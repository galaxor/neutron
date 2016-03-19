package util

import (
	"github.com/emersion/neutron/backend"
)

type EventedLabelsBackend struct {
	backend.LabelsBackend
	events backend.EventsBackend
}

func (b *EventedLabelsBackend) InsertLabel(user string, label *backend.Label) (*backend.Label, error) {
	label, err := b.LabelsBackend.InsertLabel(user, label)

	if err == nil {
		event := backend.NewLabelDeltaEvent(label.ID, backend.EventCreate, label)
		b.events.InsertEvent(user, event)
	}

	return label, err
}

func (b *EventedLabelsBackend) UpdateLabel(user string, update *backend.LabelUpdate) (*backend.Label, error) {
	label, err := b.LabelsBackend.UpdateLabel(user, update)

	if err == nil {
		event := backend.NewLabelDeltaEvent(label.ID, backend.EventUpdate, label)
		b.events.InsertEvent(user, event)
	}

	return label, err
}

func (b *EventedLabelsBackend) DeleteLabel(user, id string) error {
	err := b.LabelsBackend.DeleteLabel(user, id)

	if err == nil {
		event := backend.NewLabelDeltaEvent(id, backend.EventDelete, nil)
		b.events.InsertEvent(user, event)
	}

	return err
}

func NewEventedLabelsBackend(bkd backend.LabelsBackend, events backend.EventsBackend) backend.LabelsBackend {
	return &EventedLabelsBackend{
		LabelsBackend: bkd,
		events: events,
	}
}
