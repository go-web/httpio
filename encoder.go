package httpio

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"gopkg.in/yaml.v2"
)

// DefaultEncoders contains the default list of encoders per MIME type.
var DefaultEncoders = EncoderGroup{
	"xml":  EncoderMakerFunc(func(w http.ResponseWriter) Encoder { return &xmlEncoder{w} }),
	"json": EncoderMakerFunc(func(w http.ResponseWriter) Encoder { return &jsonEncoder{w} }),
	"yaml": EncoderMakerFunc(func(w http.ResponseWriter) Encoder { return &yamlEncoder{w} }),
}

type (
	// An Encoder encodes data from v.
	Encoder interface {
		Encode(v interface{}) error
	}

	// An EncoderGroup maps MIME types to EncoderMakers.
	EncoderGroup map[string]EncoderMaker

	// An EncoderMaker creates and returns a new Encoder.
	EncoderMaker interface {
		NewEncoder(w http.ResponseWriter) Encoder
	}

	// EncoderMakerFunc is an adapter for creating EncoderMakers
	// from functions.
	EncoderMakerFunc func(w http.ResponseWriter) Encoder
)

// NewEncoder implements the EncoderMaker interface.
func (f EncoderMakerFunc) NewEncoder(w http.ResponseWriter) Encoder {
	return f(w)
}

type xmlEncoder struct {
	w http.ResponseWriter
}

func (xe *xmlEncoder) Encode(v interface{}) error {
	xe.w.Header().Set("Content-Type", "application/xml")
	return xml.NewEncoder(xe.w).Encode(v)
}

type jsonEncoder struct {
	w http.ResponseWriter
}

func (je *jsonEncoder) Encode(v interface{}) error {
	je.w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(je.w).Encode(v)
}

type yamlEncoder struct {
	w http.ResponseWriter
}

func (ye *yamlEncoder) Encode(v interface{}) error {
	ye.w.Header().Set("Content-Type", "text/yaml")
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = ye.w.Write(b)
	return err
}
