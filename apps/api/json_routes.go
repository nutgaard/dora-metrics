package main

import (
	"encoding/xml"
	"github.com/go-chi/chi/v5"
	renderPkg "github.com/unrolled/render"
	"net/http"
)

type ExampleXml struct {
	XMLName xml.Name `xml:"example" json:"-"`
	One     string   `xml:"one,attr" json:"oneone"`
	Two     string   `xml:"two,attr" json:"twotwo"`
}

func jsonRouter() http.Handler {
	render := renderPkg.New()
	router := chi.NewRouter()

	router.Get("/json", func(writer http.ResponseWriter, request *http.Request) {
		err := render.JSON(writer, 200, map[string]string{"alive": "ok"})
		if err != nil {
			writer.WriteHeader(500)
			return
		}
	})

	router.Get("/json2", func(writer http.ResponseWriter, request *http.Request) {
		render.JSON(writer, 200, ExampleXml{One: "hello", Two: "xml"})
	})

	router.Get("/xml", func(writer http.ResponseWriter, request *http.Request) {
		render.XML(writer, 200, ExampleXml{One: "hello", Two: "xml"})
	})

	return router
}
