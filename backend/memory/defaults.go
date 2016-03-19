package memory

import (
	"github.com/emersion/neutron/backend"
)

func defaultPrivateKey() string {
	return `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v1

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
tPcWTgAvgRnzzzpGfqip2Q4BmhEQ7fyjQR4=
=4Vx7
-----END PGP PRIVATE KEY BLOCK-----`
}

func defaultPublicKey() string {
	return `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

mQENBFbmiBYBCAC5a/hMztzjfXDO1fILtIqsPeRhW6h8+ua+4X+44lTxDegZ4KD3
N3BtXgwcBQ1HwXKQY5vCx3s5lVqYrLsYRM4rE7kso7u+KqHmajqPGAmWQCRBOIiv
CJc9qYM3FUxQSRsUgJfcqSl3GeJUJ4lLlE9wCiHYJcclcxGiLRM3+TwtHEFyar47
JPNVgBfoio3MADgS2o2CnNZl7sAMM0WnR+wV8HTgM3A31TATgtyGLYjM2dld19KN
sqn6QC+iPz1dkMDAxAPUJyyn+6/ZX+OrL98UsWkoMJnsSNk7GiXuj/f+GRzJeKQV
8la9eeB7DMkJ8K4rbfizddp5QBGNX6e9WOZPABEBAAG0HW5ldXRyb24gPG5ldXRy
b25AZXhhbXBsZS5vcmc+iQE+BBMBAgAoBQJW5ogWAhsDBQkB4TOABgsJCAcDAgYV
CAIJCgsEFgIDAQIeAQIXgAAKCRDuOLAhrow4F8X2B/0TbGJTFXlXFCXU14yozm/h
mUl6IMkCrQhSVsmKV1lvJuELE3M4ZJG9doOWX4CWkCksK1HZON3ZnQaHOzzIogCZ
FVSvR8onFupI3+qnn5MF7i3j3jv4kXJT3pkiXcwuFcQqP/Fu35VKcqe6BPKooU4+
vzS9z95lSjlNi2/1aDh1VSnwhBS10moOPK0DHSUuk+tEhVwUpPaxx1pCRUKYsNn4
GhLgkrGT8LslQye/YPtiZNWdN2EfdKoJPIIjlU/Trjf/twncS6bvwRRamjqoUcqe
9Ub6Sh12hCP31UnzFgCuyChKwDDT6q4rbalpAUwO3brC9nwKN2ZV4N0TRSqu39Hy
uQENBFbmiBYBCAC9Qa6xBK16lJZDaONM2A80df8m9amDKlIUsuZBMK8S5min3kcp
gRX6kkTsnoO0Hi3DC0dJC+NfHPyYLhjTXk/IEb499tByil9YR/KX7zRhUkr6SC8L
uGbx6E1aOhUwTV+Q9j5vsLv/DJwVS9iTZn/c7l2ceRGvzP+Vveg2H6q6IkD7yemP
7t0r4usYvso7DsLBuVhOvTdACmfkC9JUEp4Nlq41kgEBfcJil6mDKnrueLTbzKrG
KwSXgNwUwABKwB6Sr9WuLaUc2H/7unLuac+ZdxDMAXjMEVmecrGeSFO3alsnAejO
hPLG3iNUWiYjfSpi4tw88YiNM01Ai1HqNZ9zABEBAAGJASUEGAECAA8FAlbmiBYC
GwwFCQHhM4AACgkQ7jiwIa6MOBcHugf/RPM6yFNz+JmBoW/OAQJl8hIH02K8LKHt
d3mW86nzk17G6w7fKxhOqg1PgGHCh0LRiK+YSX1ziYydluJ9j0jGeXspJKDb6DAb
JEw7MFqVx8Dhx2+CTBQ+NM3oLu0W5KlvUJiHoQ1ZdeenIYL6/kS/5LBcgPtOFqGi
OPkgBoHpKBZMlgSatPbhF5VkkF3PKYXinPSXQraD/TXOHr38nwNJ8Dln4tQ4+iLJ
+aapXmYZWk7LUhhR6hxFUctPVr5yjG6rYB6owSVAIQESt755dYITiQp+gw2xX4CU
kKPymjeuNF5GKhfkP1C09xZOAC+BGfPPOkZ+qKnZDgGaERDt/KNBHg==
=1lqJ
-----END PGP PUBLIC KEY BLOCK-----`
}

func (b *DomainsBackend) Populate() {
	b.domains = []*backend.Domain{
		&backend.Domain{
			ID: "domain_id",
			Name: "example.org",
		},
	}
}

func (b *Backend) Populate() {
	b.DomainsBackend.(*DomainsBackend).Populate()

	b.data = map[string]*userData{
		"user_id": &userData{
			user: &backend.User{
				ID: "user_id",
				Name: "neutron",
				DisplayName: "Neutron",
				Addresses: []*backend.Address{
					&backend.Address{
						ID: "address_id",
						DomainID: "domain_id",
						Email: "neutron@example.org",
						Send: 1,
						Receive: 1,
						Status: 1,
						Type: 1,
						Keys: []*backend.Keypair{
							&backend.Keypair{
								ID: "keypair_id",
								PublicKey: defaultPublicKey(),
								PrivateKey: defaultPrivateKey(),
							},
						},
					},
				},
			},
			password: "neutron",
			messages: []*backend.Message{
				&backend.Message{
					ID: "message_id",
					ConversationID: "conversation_id",
					AddressID: "address_id",
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
	}

	b.InsertContact("user_id", &backend.Contact{
		Name: "Myself :)",
		Email: "neutron@example.org",
	})
}
