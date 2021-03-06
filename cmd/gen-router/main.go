package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"text/template"
)

const strTmpl = `//generated by gen-router; DO NOT EDIT
package main

import "github.com/gorilla/mux"

func newRouter() *mux.Router {
	r := mux.NewRouter(){{range .}}
	r.HandleFunc("{{.Path}}", {{.Handler}}).Methods("{{.Method}}"){{end}}
	return r
}
`

type route struct {
	Method  string
	Path    string
	Handler string
}

func main() {
	inPath := flag.String("i", "", "input file path")
	outPath := flag.String("o", "", "ouput file path")
	flag.Parse()

	in, err := ioutil.ReadFile(*inPath)
	if err != nil {
		panic(err)
	}

	var routes []route
	if err := json.Unmarshal(in, &routes); err != nil {
		panic(err)
	}

	tmpl, err := template.New("router").Parse(strTmpl)
	if err != nil {
		panic(err)
	}

	bf := new(bytes.Buffer)
	if err := tmpl.Execute(bf, routes); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(*outPath, bf.Bytes(), 0644); err != nil {
		panic(err)
	}
}
