package backend

// Stores events data.
type EventsBackend interface {
	// Insert a new event.
	InsertEvent(user string, event *Event) error
	// Get the last event.
	GetLastEvent(user string) (*Event, error)
	// Get the sum of all events after a specific one.
	GetEventsAfter(user, id string) (*Event, error)
	// Delete all user's events. This happens when the user is no longer connected.
	DeleteAllEvents(user string) error
}

type Event struct {
	ID string `json:"EventID"`

	Refresh int
	Reload int
	Notices []string

	// See https://github.com/ProtonMail/WebClient/blob/master/src/app/services/event.js#L274
	Messages []*EventMessageDelta
	Conversations []*EventConversationDelta
	MessageCounts []*MessagesCount
	ConversationCounts []*MessagesCount
	Labels []*EventLabelDelta
	Contacts []*EventContactDelta
	User *User
	//Domains
	//Members
	//Organization

	UsedSpace int `json:",omitempty"`
}

type EventAction int

const (
	EventDelete EventAction = iota
	EventCreate
	EventUpdate
)

type EventDelta struct {
	ID string
	Action EventAction
}

type EventMessageDelta struct {
	EventDelta
	Message *Message
}

func NewMessageDeltaEvent(id string, action EventAction, msg *Message) *Event {
	return &Event{
		Messages: []*EventMessageDelta{
			&EventMessageDelta{
				EventDelta: EventDelta{ID: id, Action: action},
				Message: msg,
			},
		},
	}
}

type EventConversationDelta struct {
	EventDelta
	Conversation *Conversation
}

func NewConversationDeltaEvent(id string, action EventAction, conv *Conversation) *Event {
	return &Event{
		Conversations: []*EventConversationDelta{
			&EventConversationDelta{
				EventDelta: EventDelta{ID: id, Action: action},
				Conversation: conv,
			},
		},
	}
}

type EventLabelDelta struct {
	EventDelta
	Label *Label
}

func NewLabelDeltaEvent(id string, action EventAction, label *Label) *Event {
	return &Event{
		Labels: []*EventLabelDelta{
			&EventLabelDelta{
				EventDelta: EventDelta{ID: id, Action: action},
				Label: label,
			},
		},
	}
}

type EventContactDelta struct {
	EventDelta
	Contact *Contact
}

func NewContactDeltaEvent(id string, action EventAction, contact *Contact) *Event {
	return &Event{
		Contacts: []*EventContactDelta{
			&EventContactDelta{
				EventDelta: EventDelta{ID: id, Action: action},
				Contact: contact,
			},
		},
	}
}

func NewUserEvent(user *User) *Event {
	return &Event{
		User: user,
	}
}
