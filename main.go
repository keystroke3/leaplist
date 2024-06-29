package main

import (
	"html/template"
	"io"
	"relay-kiwi/stores"

	"github.com/labstack/echo/v4"
)

const (
	maxRelays = 50
	pageSize  = 10
)

type Station struct {
	Id     string
	UserId string
}

type Tag struct {
	Id        int
	Label     string
	StationId string
}

type Relay struct {
	Id          string
	Title       string
	Destination string
	Tags        []string
	Alias       string
	Note        string
	StationId   string
}

type RelayTag struct {
	RelayId   int64
	TagId     int64
	StationId string
}

type RelayGroup struct {
	Relays []stores.Relay
	Size   int64
	Next   int64
}

type User struct {
	Id          string
	Username    string
	DisplayName string
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// A custom context that extends `echo.Context`
//

func main() {
	Serve()
}
