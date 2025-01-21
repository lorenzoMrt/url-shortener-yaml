package url_shortener_yaml

import (
	yaml "gopkg.in/yaml.v3"
	"net/http"
)

type pathUrl struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err2 := parseYaml(yml)
	if err2 != nil {
		return nil, err2
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	if pathUrls == nil {
		pathUrls = make([]pathUrl, 0)
	}
	var pathMap = make(map[string]string)
	for _, v := range pathUrls {
		pathMap[v.Path] = v.Url
	}
	return pathMap
}

func parseYaml(yml []byte) ([]pathUrl, error) {
	var urls []pathUrl
	err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}
