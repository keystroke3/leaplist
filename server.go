package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"relay-kiwi/stores"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomContext struct {
	echo.Context
}

var ctx = context.Background()
var migrations = "file://stores/migration"

var db, err = stores.NewSqliteDatabase(ctx, migrations, "./stores/kiwi.sqlite")

// Defaults to `n` or zero is if the parameter is empty or cannot be converted to int.
func (c *CustomContext) IntParam(s string, n ...int) int {
	if len(n) == 0 {
		n = []int{0}
	}
	p := c.QueryParam(s)
	num, err := strconv.Atoi(p)
	if err != nil || num == 0 {
		return n[0]
	}
	return num
}

func Serve() {
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			return next(cc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/**/*.html")),
	}

	e.Renderer = renderer

	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		cc := c.(*CustomContext)
		num := cc.IntParam("samples", 10)
		return c.Render(http.StatusOK, "home.html", num)
	})
	e.GET("/relays", func(c echo.Context) error {
		cc := c.(*CustomContext)
		num := cc.IntParam("n", 10)
		fmt.Println("***** relay n = ", num)
		relays, err := db.GetStationRelays(ctx, "choochoo")
		if err != nil {
			log.Println(err)
			return c.Render(http.StatusInternalServerError, "500.html", nil)
		}
		list := RelayGroup{
			Relays: relays,
			Size:   int64(len(relays)),
			Next:   2,
		}
		return c.Render(http.StatusOK, "relays.html", list)
	})
	e.GET("/:sc", func(c echo.Context) error {
		sc := c.Param("sc")
		relay, err := db.GetRelayByID(ctx, sc)
		if err != nil {
			return c.Render(http.StatusNotFound, "404.html", nil)
		}
		return c.Redirect(http.StatusTemporaryRedirect, relay.Destination)
	})
	e.GET("/tag/:t", func(c echo.Context) error {
		tag := c.Param("t")
		tagRelays := relaysByTag(tag)
		return c.Render(http.StatusOK, "relays.html", tagRelays)
	})
	e.GET("/404", func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404.html", nil)
	})
	e.GET("/500", func(c echo.Context) error {
		return c.Render(http.StatusInternalServerError, "500.html", nil)
	})

	api := e.Group("/api")
	api.GET("/relays", func(c echo.Context) error {
		n := c.QueryParam("n")
		if n == "" {
			n = "10"
		}
		num, _ := strconv.Atoi(n)
		if num > maxRelays {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Too many relays requested, max is %v", maxRelays))
		}
		relays := nRelays(num)
		return c.JSON(http.StatusOK, relays)
	})

}
