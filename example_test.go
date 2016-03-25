package httpio_test

import (
	"encoding/xml"
	"html/template"
	"net/http"

	"github.com/go-web/httpio"
)

type Person struct {
	XMLName  xml.Name `xml:"params" json:"-" yaml:"-" schema:"-"`
	Name     string   `xml:"name" json:"name" yaml:"name" schema:"name"`
	Age      int      `xml:"age" json:"age" yaml:"age" schema:"age"`
	Location Location `xml:"location" json:"location" yaml:"location" schema:"location"`
}

type Location struct {
	Address string `xml:"address" json:"address" yaml:"address" schema:"address"`
}

type Response struct {
	OK     bool
	Person Person
}

var ResponseT = template.Must(template.New("resp").Parse(`<h1>OK: {{.OK}}</h1>
<pre>
name: {{.Person.Name}}
age: {{.Person.Age}}
address: {{.Person.Location.Address}}
</pre>
`))

func Example() {
	// curl -F name=Bob -F age=22 -F location.address=internets localhost:8080
	// curl -H 'Content-Type: application/json' -d'{"name":"Bob","age":22,"location":{"address":"internets"}}' localhost:8080
	// Add the Accept header to get different responses.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var p Person
		_, err := httpio.Bind(r, &p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// ...
		resp := &Response{OK: true, Person: p}
		httpio.Write(w, r, resp, func() error {
			return ResponseT.Execute(w, resp)
		})
	})
	http.ListenAndServe(":8080", nil)
}
