package main

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	maxLeaps = 50
	pageSize = 10
)

type Leap struct {
	Name        string
	Original    string
	Tags        []string
	Shortcut    string
	Description string
}

type LeapList struct {
	Leaps []Leap
	Size  int
	Next  int
}

var samples = sampleLeaps(maxLeaps)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// A custom context that extends `echo.Context`
type CustomContext struct {
	echo.Context
}

//	Attempts to get the url query parameter `s` and convert it to an int.
//
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

// selects a string from a given slice of strings `choices`
func randomSelect(choices []string) string {
	n := rand.Intn(len(choices))
	return choices[n]
}

func randomUrl() string {
	return fmt.Sprintf("https://%v.%v/%v", randomSelect(nouns), randomSelect(tlds), randomSelect(nouns))
}

// provides a slice of `Leap` for testing
func sampleLeaps(n int) map[string]Leap {
	randTagsN := func(j int) []string {
		var tags []string
		for i := 0; i < j; i++ {
			tags = append(tags, randomSelect(adjectives))
		}
		return tags
	}

	l := make(map[string]Leap)
	for i := 0; i < n; i++ {
		sc := randomSelect(nouns)
		_, set := l[sc]
		if set {
			sc = fmt.Sprintf("%v%v", sc, i)
		}
		l[sc] = Leap{
			Name:        randomSelect(nouns),
			Shortcut:    sc,
			Original:    randomUrl(),
			Tags:        randTagsN(rand.Intn(6)),
			Description: randomSelect(descriptions),
		}
	}
	return l
}

func nLeaps(n int) []Leap {
	l := make([]Leap, n)
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

func leapsByTag(tag string) []Leap {
	leaps := []Leap{}
	for _, v := range samples {
		for _, t := range v.Tags {
			if t == tag {
				leaps = append(leaps, v)
				break
			}
		}
	}
	return leaps
}

func main() {
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
	e.GET("/leaps", func(c echo.Context) error {
		cc := c.(*CustomContext)
		num := cc.IntParam("n", 10)
		fmt.Println("***** leap n = ", num)
		leaps := nLeaps(num)
		list := LeapList{
			Leaps: leaps,
			Size:  len(leaps),
			Next:  2,
		}
		return c.Render(http.StatusOK, "leaps.html", list)
	})
	e.GET("/:sc", func(c echo.Context) error {
		sc := c.Param("sc")
		sample, set := samples[sc]
		if !set {
			return c.Render(http.StatusNotFound, "404.html", nil)
		}
		return c.Redirect(http.StatusTemporaryRedirect, sample.Original)
	})
	e.GET("/tag/:t", func(c echo.Context) error {
		tag := c.Param("t")
		tagLeaps := leapsByTag(tag)
		return c.Render(http.StatusOK, "leaps.html", tagLeaps)
	})
	e.GET("/404", func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404.html", nil)
	})
	e.GET("/500", func(c echo.Context) error {
		return c.Render(http.StatusInternalServerError, "500.html", nil)
	})

	api := e.Group("/api")
	api.GET("/leaps", func(c echo.Context) error {
		n := c.QueryParam("n")
		if n == "" {
			n = "10"
		}
		num, _ := strconv.Atoi(n)
		if num > maxLeaps {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Too many leaps requested, max is %v", maxLeaps))
		}
		leaps := nLeaps(num)
		return c.JSON(http.StatusOK, leaps)
	})

	e.Logger.Fatal(e.Start(":7000"))
}
