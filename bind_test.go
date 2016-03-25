package httpio

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type Params struct {
	XMLName  xml.Name `xml:"params" json:"-" yaml:"-" schema:"-"`
	Name     string   `xml:"name" json:"name" yaml:"name" schema:"name"`
	Age      int      `xml:"age" json:"age" yaml:"age" schema:"age"`
	Location struct {
		Address string `xml:"address" json:"address" yaml:"address" schema:"address"`
	} `xml:"location" json:"location" yaml:"location" schema:"location"`
}

func TestBind(t *testing.T) {
	cases := []struct {
		Request *http.Request
		Code    int
	}{
		{
			Request: &http.Request{
				URL: &url.URL{Path: "/"},
				PostForm: url.Values{
					"name":             {"Bob"},
					"age":              {"13"},
					"location.address": {"internets"},
				},
				Header: http.Header{"Content-Type": {"form-urlencoded"}},
			},
			Code: http.StatusOK,
		},
		{
			Request: &http.Request{
				URL:    &url.URL{Path: "/"},
				Header: http.Header{"Content-Type": {"application/xml"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`
				<params>
					<name>Bob</name>
					<age>13</age>
					<location>
						<address>internets</address>
					</location>
				</params>
				`))),
			},
			Code: http.StatusOK,
		},
		{
			Request: &http.Request{
				URL:    &url.URL{Path: "/"},
				Header: http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`
				{
					"name":"Bob",
					"age":13,
					"location":{
						"address":"internets"
					}
				}
				`))),
			},
			Code: http.StatusOK,
		},
		{
			Request: &http.Request{
				URL:    &url.URL{Path: "/"},
				Header: http.Header{"Content-Type": {"text/yaml"}},
				Body: ioutil.NopCloser(bytes.NewReader(
					[]byte("name: Bob\nage: 13\nlocation:\n address: internets\n"),
				)),
			},
			Code: http.StatusOK,
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var in Params
		_, err := Bind(r, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch {
		case in.Name != "Bob", in.Age != 13, in.Location.Address != "internets":
			m := fmt.Sprintf("invalid form: %#v", in)
			http.Error(w, m, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	for i, tc := range cases {
		var b bytes.Buffer
		w := &httptest.ResponseRecorder{Body: &b}
		mux.ServeHTTP(w, tc.Request)
		if w.Code != tc.Code {
			t.Errorf("test %d failed: want code %d, have %d\n%s",
				i, tc.Code, w.Code, b.Bytes())
		}
	}
}

func TestWrite(t *testing.T) {
	data := &Params{
		Name: "Bob",
		Age:  13,
	}
	data.Location.Address = "internets"
	cases := []struct {
		Accept string
		Error  error
		Data   string
	}{
		{
			Accept: "application/xml",
			Error:  nil,
			Data:   "<params><name>Bob</name><age>13</age><location><address>internets</address></location></params>",
		},
		{
			Accept: "application/json",
			Error:  nil,
			Data:   `{"name":"Bob","age":13,"location":{"address":"internets"}}` + "\n",
		},
		{
			Accept: "text/yaml",
			Error:  nil,
			Data:   "name: Bob\nage: 13\nlocation:\n  address: internets\n",
		},
		{
			Accept: "use-default-encoder",
			Error:  nil,
			Data:   `{"name":"Bob","age":13,"location":{"address":"internets"}}` + "\n",
		},
	}
	for i, tc := range cases {
		var b bytes.Buffer
		w := &httptest.ResponseRecorder{Body: &b}
		r := &http.Request{
			Header: http.Header{"Accept": {tc.Accept}},
		}
		err := Write(w, r, data, func() error {
			return json.NewEncoder(w).Encode(data)
		})
		if err != tc.Error {
			t.Errorf("test %d failed to encode %s: want %v, have %v",
				i, tc.Accept, tc.Error, err)
		}
		have := b.String()
		if have != tc.Data {
			t.Errorf("test %d failed:\nwant %q\nhave %q",
				i, tc.Data, have)
		}
	}
}
