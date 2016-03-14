package backend

type Label struct {
	ID string
	Name string
	Color string
	Display int
	Order int
}

func GetLabels(id string) (labels []*Label, err error) {
	labels = []*Label{
		&Label{
			ID: "label_id",
			Name: "Hey!",
			Color: "#7272a7",
			Display: 1,
			Order: 1,
		},
	}
	return
}
