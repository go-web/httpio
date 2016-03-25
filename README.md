# httpio

[![GoDoc](https://godoc.org/github.com/go-web/httpio?status.svg)](http://godoc.org/github.com/go-web/httpio)

httpio provides input decoders and output encoders. The idea is that you
define a struct for the input parameters and httpio will parse it based
on the Content-Type. It supports XML, JSON, YAML, Form and Multipart Form.

For the output encoders, same thing. You define a struct and httpio writes
the response based on the Accept header. If no encoder is found (in case
Accept is not) you can provide a default encoder of your own.

Install:

```
go get github.com/go-web/httpio
```

See the [example](example_test.go) for details.
