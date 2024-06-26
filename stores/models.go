package stores

type Relay struct {
	Id          int64
	Title       string
	Alias       string
	Destination string
	Target      string
	Note        string
	Description string
}

type User struct {
	Id          string
	Username    string
	DisplayName string
}
type Tag struct {
	Id    int64
	Label string
}

type Station struct {
	Id     string
	UserId string
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
