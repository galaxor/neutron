package backend

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
	//User
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

type EventContactDelta struct {
	EventDelta
	Contact *Contact
}
