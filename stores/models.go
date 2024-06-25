package stores

type Relay struct {
	Name        string
	Original    string
	Value       string
	Description string
}
type Tag struct {
	Name string
}

type RelayTags struct {
	RelayId string
	TagId   string
}

type SqlRelay struct {
	Name        string
	Original    string
	Value       string
	Description string
}
