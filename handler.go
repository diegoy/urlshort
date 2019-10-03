package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
		redirectURL, urlPresent := pathsToUrls[req.URL.Path]
		fmt.Printf("Received request: %s\n", req.URL.Path)

		if urlPresent {
			http.Redirect(responseWriter, req, redirectURL, 302)
		} else {
			fallback.ServeHTTP(responseWriter, req)
		}
	})
}

type tuple struct {
	Path string
	URL  string
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urls := []tuple{}

	err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return nil, err
	}
	fmt.Printf("---- yaml \n%v\n\n", urls)

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
		fmt.Printf("Received request: %s\n", req.URL.Path)

		foundURL := ""
		for _, tuple := range urls {
			if tuple.Path == req.URL.Path {
				foundURL = tuple.URL
				break
			}
		}

		if foundURL != "" {
			http.Redirect(responseWriter, req, foundURL, 302)
		} else {
			fallback.ServeHTTP(responseWriter, req)
		}
	}), nil
}

// JSONHandler will handle urls based on the provided JSON file
// Expected Format:
//
// [
// 	{ path: '/some-path', url: 'https://example.com'}
// ]
//
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urls := []tuple{}

	err := json.Unmarshal(jsonBytes, &urls)
	if err != nil {
		return nil, err
	}
	fmt.Printf("---- json \n%v\n\n", urls)

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
		fmt.Printf("Received request: %s\n", req.URL.Path)

		foundURL := ""
		for _, tuple := range urls {
			if tuple.Path == req.URL.Path {
				foundURL = tuple.URL
				break
			}
		}

		if foundURL != "" {
			http.Redirect(responseWriter, req, foundURL, 302)
		} else {
			fallback.ServeHTTP(responseWriter, req)
		}
	}), nil
}
