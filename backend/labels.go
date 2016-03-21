package backend

// Stores labels data.
type LabelsBackend interface {
	// List all user's labels.
	ListLabels(user string) ([]*Label, error)
	// Insert a new label.
	InsertLabel(user string, label *Label) (*Label, error)
	// Update an existing label.
	UpdateLabel(user string, update *LabelUpdate) (*Label, error)
	// Delete a label.
	DeleteLabel(user, id string) error
}

// A message label.
type Label struct {
	ID string
	Name string
	Color string
	Display int
	Order int
}

// System labels
const (
	InboxLabel string = "0"
	DraftsLabel = "1"
	SentLabel = "2"
	TrashLabel = "3"
	SpamLabel = "4"
	ArchiveLabel = "6"
	StarredLabel = "10"
)

// A request to update a label.
// Fields set to true will be updated with values in Label.
type LabelUpdate struct {
	Label *Label
	Name bool
	Color bool
	Display bool
	Order bool
}

// Apply this update on a label.
func (update *LabelUpdate) Apply(label *Label) {
	updated := update.Label

	if updated.ID != label.ID {
		panic("Cannot apply update on a label with a different ID")
	}

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
}
