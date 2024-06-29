package stores

type Relay struct {
	Id          int64
	Title       string
	Alias       string
	Destination string
	Target      string
	Note        string
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
	Title  string
	Domain string
	UserId string
}

type RelayTags struct {
	RelayId string
	TagId   string
}
