package httpio

import (
	"errors"
	"net/http"
	"strings"
)

// ErrInvalidFormat indicates the data passed to encode or decode is invalid.
var ErrInvalidFormat = errors.New("invalid data format")

// Bind binds the request input to v based on the request's content type.
// v must be a pointer to struct.
//
// The struct it points to will store request input using DefaultDecoders,
// that is, one of XML, JSON, YAML, Form, or Multipart Form. Accepted tags
// in the given struct are xml, json, yaml and schema.
//
// XML and JSON are parsed using encoding/xml and encoding/json respectively,
// while YAML uses gopkg.in/yaml.v2 and Forms use github.com/gorilla/schema.
func Bind(r *http.Request, v interface{}) (mime string, err error) {
	t := r.Header.Get("Content-Type")
	for k, f := range DefaultDecoders {
		if strings.Contains(t, k) {
			return k, f.NewDecoder(r).Decode(v)
		}
	}
	return t, ErrInvalidFormat
}

// Write looks at the acceptable formats in the request's accept header
// and encodes v into w using DefaultEncoders. The optional function f
// is called when no format match the list of default encoders, one of
// XML, JSON or YAML.
func Write(w http.ResponseWriter, r *http.Request, v interface{}, f func() error) error {
	t := r.Header.Get("Accept")
	for k, f := range DefaultEncoders {
		if strings.Contains(t, k) {
			return f.NewEncoder(w).Encode(v)
		}
	}
	if f != nil {
		return f()
	}
	return ErrInvalidFormat
}
