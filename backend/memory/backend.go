package memory

import (
	"github.com/emersion/neutron/backend"
)

type Backend struct {
	data map[string]*userData
}

type userData struct {
	user *backend.User
	password string
	contacts []*backend.Contact
	conversations []*backend.Conversation
	messages []*backend.Message
	labels []*backend.Label
}

func New() backend.Backend {
	return &Backend{
		data: map[string]*userData{
			"user_id": &userData{
				user: &backend.User{
					ID: "user_id",
					Name: "neutron",
					DisplayName: "Neutron",
					PublicKey: defaultPublicKey(),
					EncPrivateKey: defaultPrivateKey(),
					Addresses: []*backend.Address{
						&backend.Address{
							ID: "address_id",
							DomainID: "domain_id",
							Email: "neutron@example.org",
							Send: 1,
							Receive: 1,
						},
					},
				},
				password: "neutron",
				contacts: []*backend.Contact{
					&backend.Contact{
						ID: "contact_id",
						Name: "Myself :)",
						Email: "neutron@example.org",
					},
				},
				conversations: []*backend.Conversation{
					&backend.Conversation{
						ID: "conversation_id",
						Subject: "Hi from Neutron!",
						Senders: []*backend.Email{ &backend.Email{"neutron@example.org", "Neutron"} },
						Recipients: []*backend.Email{ &backend.Email{"neutron@example.org", "Neutron"} },
						Time: 1458073557,
					},
				},
				messages: []*backend.Message{
					&backend.Message{
						ID: "message_id",
						ConversationID: "conversation_id",
						Subject: "Hello World",
						Sender: &backend.Email{"neutron@example.org", "Neutron"},
						ToList: []*backend.Email{ &backend.Email{"neutron@example.org", "Neutron"} },
						Time: 1458073557,
						Body: "Hey! How are you today?",
						LabelIDs: []string{"0"},
					},
				},
				labels: []*backend.Label{
					&backend.Label{
						ID: "label_id",
						Name: "Hey!",
						Color: "#7272a7",
						Display: 1,
						Order: 1,
					},
				},
			},
		},
	}
}
