package httpio

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/schema"
	"gopkg.in/yaml.v2"
)

// DefaultDecoders contains the default list of decoders per MIME type.
var DefaultDecoders = DecoderGroup{
	"xml":                 DecoderMakerFunc(func(r *http.Request) Decoder { return xml.NewDecoder(r.Body) }),
	"json":                DecoderMakerFunc(func(r *http.Request) Decoder { return json.NewDecoder(r.Body) }),
	"yaml":                DecoderMakerFunc(func(r *http.Request) Decoder { return &yamlDecoder{r.Body} }),
	"form-urlencoded":     DecoderMakerFunc(func(r *http.Request) Decoder { return &formDecoder{r, false} }),
	"multipart/form-data": DecoderMakerFunc(func(r *http.Request) Decoder { return &formDecoder{r, true} }),
}

type (
	// A Decoder decodes data into v.
	Decoder interface {
		Decode(v interface{}) error
	}

	// A DecoderGroup maps MIME types to DecoderMakers.
	DecoderGroup map[string]DecoderMaker

	// A DecoderMaker creates and returns a new Decoder.
	DecoderMaker interface {
		NewDecoder(r *http.Request) Decoder
	}

	// DecoderMakerFunc is an adapter for creating DecoderMakers
	// from functions.
	DecoderMakerFunc func(r *http.Request) Decoder
)

// NewDecoder implements the DecoderMaker interface.
func (f DecoderMakerFunc) NewDecoder(r *http.Request) Decoder {
	return f(r)
}

type yamlDecoder struct {
	r io.Reader
}

func (yd *yamlDecoder) Decode(v interface{}) error {
	b, err := ioutil.ReadAll(yd.r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

type formDecoder struct {
	r         *http.Request
	multipart bool
}

func (fd *formDecoder) Decode(v interface{}) error {
	var f func() error
	if fd.multipart {
		f = func() error { return fd.r.ParseMultipartForm(1 << 12) }
	} else {
		f = func() error { return fd.r.ParseForm() }
	}
	if err := f(); err != nil {
		return err
	}
	return schema.NewDecoder().Decode(v, fd.r.Form)
}
