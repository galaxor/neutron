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
	//Labels
	//Contacts
	//User
	//Domains
	//Members
	//Organization

	UsedSpace int `json:omitempty`
}

type EventAction int

const (
	EventDelete EventAction = 0
	EventCreate = 1
	EventUpdate = 2
)

type EventDelta struct {
	ID string
	Action EventAction
}

type EventMessageDelta struct {
	EventDelta
	Message *Message
}

type EventConversationDelta struct {
	EventDelta
	Conversation *Conversation
}
