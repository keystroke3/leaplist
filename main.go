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

type Leap struct {
	Name     string
	Original string
	Tags     []string
	Shortcut string
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var tlds = []string{"com", "net", "org", "gov", "edu", "info", "biz", "co", "io", "pro"}
var lorem = []string{
	"animi", "voluptatem", "repellendusg", "laboriosam", "facilis",
	"assumenda", "autem", "vero", "veniam", "repudiandae", "repellat", "corporis",
	"ullam", "temporibus", "incidunt", "harum", "nulla", "Animi", "praesentium", "quasi", "ipsam",
}

// selects a string from a given slice of strings `choices`
func randomSelect(choices []string) string {
	n := rand.Intn(len(choices))
	return choices[n]
}

func randomUrl() string {
	return fmt.Sprintf("https://%v.%v/%v", randomSelect(lorem), randomSelect(tlds), randomSelect(lorem))
}

// provides a slice of `Leap` for testing
func sampleLeaps(n int) map[string]Leap {

	randTagsN := func(j int) []string {
		var tags []string
		for i := 0; i < j; i++ {
			tags = append(tags, randomSelect(lorem))
		}
		return tags
	}

	l := make(map[string]Leap)
	for i := 0; i < n; i++ {
		sc := randomSelect(lorem)
		l[sc] = Leap{
			Name:     randomSelect(lorem),
			Shortcut: sc,
			Original: randomUrl(),
			Tags:     randTagsN(rand.Intn(6)),
		}
	}
	return l
}

var samples = sampleLeaps(10)

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/**/*.html")),
	}

	e.Renderer = renderer

	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		n := c.QueryParam("samples")
		if n == "" {
			n = "5"
		}
		num, _ := strconv.Atoi(n)
		return c.Render(http.StatusOK, "home.html", num)
	})

	e.GET("/leaps", func(c echo.Context) error {
		n := c.QueryParam("n")
		if n == "" {
			n = "5"
		}
		num, _ := strconv.Atoi(n)
		leaps := nLeaps(num)
		fmt.Println("** samples", len(leaps))
		return c.Render(http.StatusOK, "leaps.html", leaps)
	})
	e.GET("/:sc", func(c echo.Context) error {
		sc := c.Param("sc")
		return c.Redirect(http.StatusTemporaryRedirect, samples[sc].Original)
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

	e.GET("/api/leaps", func(c echo.Context) error {
		n := c.QueryParam("n")
		if n == "" {
			n = "5"
		}
		num, _ := strconv.Atoi(n)

		leaps := sampleLeaps(num)
		fmt.Println("** samples", len(leaps))
		return c.JSON(http.StatusOK, leaps)
	})

	e.Logger.Fatal(e.Start(":5000"))
}
