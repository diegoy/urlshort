package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gophercises/urlshort"
)

func main() {
	fileName := flag.String("file", "default", "file name")
	format := flag.String("format", "yaml", "type of the file: [yaml, json]")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	if strings.Compare(*format, "yaml") == 0 {
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		yaml := getUrlsYaml(*fileName)
		yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	} else if strings.Compare(*format, "json") == 0 {
		json := getUrlsJSON(*fileName)
		jsonHandler, err := urlshort.JSONHandler([]byte(json), mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
	} else {
		panic("Unknown file format")
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func getUrlsYaml(yamlFile string) string {
	if yamlFile != "default" {
		return readFile(yamlFile)
	}

	return `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
}

func getUrlsJSON(fileName string) string {
	if fileName != "default" {
		return readFile(fileName)
	}

	return `
[ {"path": "/test", "url": "https://www.uol.com.br"} ]
	`
}

func readFile(fileName string) string {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Printf(`Oops file "%s" not found`, fileName)
		panic(err)
	}

	return string(file)
}
