package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func main() {
	SampleUse()
}

const HTML = `
<html>
	<head></head>
	<body>
		{{if .Visitable}}
		<p>if block</p>
		{{end}}
		
		<p>server: {{.Name}} {{.Location}}</p>

		<p>printf: {{printf "%v" .}}</p>
	</body>
</html>
`

var templates = template.Must(template.New("sample.html").Parse(HTML))

type Server struct {
	Name      string
	Location  string
	Visitable bool
}

func SampleUse() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &Server{"Name", "shanghai", true}
		err := templates.ExecuteTemplate(w, "sample.html", data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}))
	http.ListenAndServe(":8080", nil)
}
