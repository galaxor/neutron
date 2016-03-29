package memory

import (
	"github.com/emersion/neutron/backend"
)

/*func (b *Domains) Populate() {
	b.domains = []*backend.Domain{
		&backend.Domain{
			ID: "domain_id",
			Name: "example.org",
		},
	}
}*/

func Populate(b *backend.Backend) (err error) {
	// TODO
	//b.DomainsBackend.(*Domains).Populate()

	user, err := b.InsertUser(&backend.User{
		Name: "neutron",
		DisplayName: "Neutron",
		Addresses: []*backend.Address{
			&backend.Address{
				DomainID: "domain_id",
				Email: "neutron@example.org",
				Send: 1,
				Receive: 1,
				Status: 1,
				Type: 1,
			},
		},
	}, "neutron")
	if err != nil {
		return
	}

	b.InsertContact(user.ID, &backend.Contact{
		Name: "Myself :)",
		Email: "neutron@example.org",
	})

	b.InsertLabel(user.ID, &backend.Label{
		Name: "Hey!",
		Color: "#7272a7",
		Display: 1,
		Order: 1,
	})

	b.InsertMessage(user.ID, &backend.Message{
		ID: "message_id",
		ConversationID: "conversation_id",
		AddressID: "address_id",
		Subject: "Hello World",
		Sender: &backend.Email{"neutron@example.org", "Neutron"},
		ToList: []*backend.Email{ &backend.Email{"neutron@example.org", "Neutron"} },
		Time: 1458073557,
		Body: "Hey! How are you today?",
		LabelIDs: []string{"0"},
	})

	return
}
