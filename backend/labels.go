package backend

const (
	InboxLabel string = "0"
	DraftsLabel = "1"
	SentLabel = "2"
	TrashLabel = "3"
	SpamLabel = "4"
	ArchiveLabel = "6"
	StarredLabel = "10"
)

type Label struct {
	ID string
	Name string
	Color string
	Display int
	Order int
}

type LabelUpdate struct {
	Label *Label
	Name bool
	Color bool
	Display bool
	Order bool
}
