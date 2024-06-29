package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"relay-kiwi/stores"
)

var samples = sampleRelays(maxRelays)
var user = User{Id: "foouser", Username: "foo", DisplayName: "foo"}
var station = Station{Id: "choochoo", UserId: user.Id}

// selects a string from a given slice of strings `choices`
func randomSelect(choices []string) string {
	n := rand.Intn(len(choices))
	return choices[n]
}

func randomUrl() string {
	return fmt.Sprintf("https://%v.%v/%v", randomSelect(nouns), randomSelect(tlds), randomSelect(nouns))
}

// provides a slice of `relay` for testing
func sampleRelays(n int, stationID ...string) map[string]Relay {
	sID := station.Id
	if len(stationID) > 0 {
		sID = stationID[0]
	}

	randTagsN := func(j int) []string {
		var tags []string
		for i := 0; i < j; i++ {
			tags = append(tags, randomSelect(adjectives))
		}
		return tags
	}

	l := make(map[string]Relay)
	for i := 0; i < n; i++ {
		sc := randomSelect(nouns)
		_, set := l[sc]
		if set {
			sc = fmt.Sprintf("%v%v", sc, i)
		}
		l[sc] = Relay{
			Title:       randomSelect(nouns),
			Alias:       sc,
			Destination: randomUrl(),
			Tags:        randTagsN(rand.Intn(6)),
			Note:        randomSelect(notes),
			StationId:   sID,
		}
	}
	return l
}

func nRelays(n int) []Relay {
	l := make([]Relay, n)
	i := 0
	for _, v := range samples {
		if i == n {
			break
		}
		l[i] = v
		i++
	}
	return l
}

func relaysByTag(tag string) []Relay {
	relays := []Relay{}
	for _, v := range samples {
		for _, t := range v.Tags {
			if t == tag {
				relays = append(relays, v)
				break
			}
		}
	}
	return relays
}

func populateTestDB() {
	ctx := context.Background()
	migrations := "file://stores/migration"
	db, err := stores.NewSqliteDatabase(ctx, migrations, "./stores/kiwi.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.CreateUser(ctx, user.Username, user.DisplayName, "fooPass")
	if err != nil {
		log.Fatal(err)
	}
	err = db.CreateStation(ctx, station.Id, user.Id)
	if err != nil {
		log.Fatal(err)
	}
	relays := sampleRelays(20)
	for _, r := range relays {
		db.CreateRelay(ctx, r.Title, r.Alias, r.Destination, r.Note, r.StationId)
	}

}

var tlds = []string{"com", "net", "org", "gov", "edu", "info", "biz", "co", "io", "pro"}

var nouns = []string{
	"electronegative", "towny", "miniaturist", "dermopteran", "ostosis", "hierarchism", "roaming", "bakeware", "ciré", "undersurface",
	"landskip", "zariba", "concurrent", "metalliding", "megalosaur", "bladderworm", "plastogene", "psammite", "presternum", "polysemy",
	"footy", "opioid", "macchiato", "ferrochromium", "weedkiller", "anglepod", "sulu", "eudiometer", "calicoback", "sampling",
	"cratering", "ap", "korat", "quantile", "watercooler", "swingby", "til", "bimolecular", "nebuchadnezzar", "bushelbasket",
	"nonbeing", "skag", "shammer", "ashes", "paratransit", "theelol", "matchmark", "covin", "intervenient", "thioaldehyde",
}

var adjectives = []string{
	"photothermic", "vagarious", "sovereigntist", "tuberose", "wire", "tantalous", "apod", "wunderbar", "cotemporary", "interradial",
	"honeycomb", "missile", "countertenor", "chapleted", "dullsville", "addicting", "ripstop", "shaly", "stroganoff", "textured",
	"indivertible", "cash", "pocky", "overjoyed", "lowering", "chaotic", "pea", "an", "windless", "noteless",
	"pilular", "cutty", "uncoachable", "knightly", "implied", "curve", "geomagnetic", "pithecanthropine", "salaried", "warmhearted",
	"ellipsoid", "deliverable", "maple", "yauld", "tent", "burdened", "shorthand", "antilithic", "frustrate", "leftish",
}

var notes = []string{
	"of or involving both light and heat",
	"wandering or roaming",
	"of or expressing support for making the province of Quebec essentially independent from Canada",
	"atuberousa",
	"made of wire or wirework",
	"of, derived from, or containing tantalum, esp trivalent tantalum",
	"aapodala",
	"wonderful",
	"acontemporarya",
	"situated between rays or radii",
	"of, like, or patterned after a honeycomb",
	"throwing or shooting missiles",
	"of, for, or having the range of a countertenor",
	"wearing a wreath or garland on the head",
	"very dull, boring, tedious, etc",
	"relating to or causing addiction or addictive",
	"designating or of a fabric, esp nylon, woven with extra threads in a pattern to make runs or tears less likely and used for parachutes, garments, etc",
	"of, like, or containing shale",
	"cooked with sour cream, onions, mushrooms, etc",
	"having a particular kind of texture, esp one that is uneven, not smooth, easily perceived by touching, etc",
	"that cannot be diverted or turned aside",
	"of, for, requiring, or made with cash",
	"of or having the pox",
	"feeling great joy",
	"dark, as if about to rain or snow or overcast",
	"of or having to do with the theories, dynamics, etc of mathematical chaos",
	"designating a family Fabaceae, order Fabales of leguminous, dicotyledonous plants, including peanuts, clover, vetch, alfalfa, and many beans",
	"any one",
	"out of breath",
	"unmusical",
	"of or like a pill or pills",
	"short",
	"not coachable or specif, not responsive to coaching, as because of temperament, stubbornness, etc",
	"consisting of knights",
	"involved, suggested, or understood without being openly or directly expressed",
	"curved",
	"of or pertaining to the magnetic properties of the earth",
	"of, belonging to, or resembling a former genus iPithecanthropusi, now classified as iHomo erectusi of extinct early humans, who lived in Java, China, Europe, and Africa",
	"receiving or yielding a salary",
	"kind, sympathetic, friendly, loving, etc",
	"of or shaped like an ellipsoid",
	"capable of being delivered",
	"flavored with maple",
	"active, nimble, etc",
	"of or like a tent",
	"designating the vessel responsible for taking action to avoid colliding with another vessel",
	"using or written in shorthand",
	"preventing the formation or development of calculi, as of the urinary tract",
	"frustrated or baffled or defeated",
	"inclined to be leftist or a leftist",
}
