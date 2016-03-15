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
					NotificationEmail: "neutron@example.org",
				},
				password: "neutron",
				contacts: []*backend.Contact{
					&backend.Contact{
						ID: "contact_id",
						Name: "Myself :)",
						Email: "neutron@example.org",
					},
				},
				conversations: []*backend.Conversation{},
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
